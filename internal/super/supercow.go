package super

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	"github.com/Code-Hex/Neo-cowsay/internal/screen"
	runewidth "github.com/mattn/go-runewidth"
)

func getNoSaidCow(opts ...cowsay.Option) (string, error) {
	opts = append(opts, cowsay.Thoughts(' '))
	cow, err := cowsay.New(opts...)
	if err != nil {
		return "", err
	}
	return cow.GetCow()
}

// RunSuperCow runs super cow mode animation on the your terminal
func RunSuperCow(phrase string, opts ...cowsay.Option) error {
	cow, err := cowsay.New(opts...)
	if err != nil {
		return err
	}
	balloon := cow.Balloon(phrase)
	blank := createBlankSpace(balloon)

	said, err := cow.GetCow()
	if err != nil {
		return err
	}

	notSaid, err := getNoSaidCow(opts...)
	if err != nil {
		return err
	}

	saidCow := balloon + said
	saidCowLines := strings.Count(saidCow, "\n") + 1

	// When it is higher than the height of the terminal
	h := screen.Height()
	if saidCowLines > h {
		return errors.New("too height messages")
	}

	notSaidCowLines := strings.Split(blank+notSaid, "\n")

	renderer := newRenderer(saidCow, notSaidCowLines)

	screen.SaveState()
	screen.HideCursor()
	screen.Clear()

	go renderer.createFrames(cow)

	renderer.render()

	screen.UnHideCursor()
	screen.RestoreState()

	return nil
}

func createBlankSpace(balloon string) string {
	var buf strings.Builder
	l := strings.Count(balloon, "\n")
	for i := 0; i < l; i++ {
		buf.WriteRune('\n')
	}
	return buf.String()
}

func maxLen(cow []string) int {
	max := 0
	for _, line := range cow {
		l := runewidth.StringWidth(line)
		if max < l {
			max = l
		}
	}
	return max
}

type renderer struct {
	max         int
	middle      int
	screenWidth int
	heightDiff  int
	frames      chan string

	saidCow         string
	notSaidCowLines []string

	quit chan os.Signal
}

func newRenderer(saidCow string, notSaidCowLines []string) *renderer {
	w, cowsWidth := screen.Width(), maxLen(notSaidCowLines)
	max := w + cowsWidth

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return &renderer{
		max:             max,
		middle:          max / 2,
		screenWidth:     w,
		heightDiff:      screen.Height() - strings.Count(saidCow, "\n") + 1,
		frames:          make(chan string, max),
		saidCow:         saidCow,
		notSaidCowLines: notSaidCowLines,
		quit:            quit,
	}
}

const (
	// Frequency the color changes
	magic = 2

	span    = 30 * time.Millisecond
	standup = 3 * time.Second
)

func (r *renderer) createFrames(cow *cowsay.Cow) {
	const times = standup / span
	var buf strings.Builder
	for x, i := 0, 0; i <= r.max; i++ {
		if i == r.middle {
			posx := r.posX(i)
			for k := 0; k < int(times); k++ {
				base := x * 70
				// draw colored cow
				screen.MoveTo(cow.Aurora(base, r.saidCow), posx, r.posY(0))
				r.frames <- screen.Flush()
				if k%magic == 0 {
					x++
				}
			}
		} else {
			var n int
			if i > r.screenWidth {
				n = i - r.screenWidth
			}

			base := x * 70
			buf.Reset()
			for _, line := range r.notSaidCowLines {
				if i > r.screenWidth {
					if n < len(line) {
						buf.WriteString(line[n:])
					}
				} else if i > len(line) {
					buf.WriteString(line)
				} else {
					buf.WriteString(line[:i])
				}
				buf.WriteString("\n")
			}
			screen.MoveTo(cow.Aurora(base, buf.String()), r.posX(i), r.posY(0))
			r.frames <- screen.Flush()
		}
		if i%magic == 0 {
			x++
		}
	}
	close(r.frames)
}

func (r *renderer) render() {
	initCh := make(chan struct{}, 1)
	initCh <- struct{}{}

	for view := range r.frames {
		select {
		case <-r.quit:
			screen.Clear()
			return
		case <-initCh:
		case <-time.After(span):
		}
		io.Copy(screen.Stdout, strings.NewReader(view))
	}
}

func (r *renderer) posX(i int) int {
	posx := r.screenWidth - i
	if posx < 1 {
		posx = 1
	}
	return posx
}

func (r *renderer) posY(i int) int {
	return r.heightDiff + i - 1
}
