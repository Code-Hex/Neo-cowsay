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

	saidCowLines := strings.Split(balloon+said, "\n")

	// When it is higher than the height of the terminal
	h := screen.Height()
	if len(saidCowLines) > h {
		return errors.New("too height messages")
	}

	notSaidCowLines := strings.Split(blank+notSaid, "\n")

	renderer := newRenderer(saidCowLines, notSaidCowLines)

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

	saidCowLines    []string
	notSaidCowLines []string

	quit chan os.Signal
}

func newRenderer(saidCowLines, notSaidCowLines []string) *renderer {
	w, cowsWidth := screen.Width(), maxLen(notSaidCowLines)
	max := w + cowsWidth

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return &renderer{
		max:             max,
		middle:          max / 2,
		screenWidth:     w,
		heightDiff:      screen.Height() - len(saidCowLines),
		frames:          make(chan string, max),
		saidCowLines:    saidCowLines,
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
	for x, i := 0, 0; i <= r.max; i++ {
		if i == r.middle {
			posx := r.posX(i)
			for k := 0; k < int(times); k++ {
				base := x * 70
				// draw colored cow
				for j, line := range r.saidCowLines {
					screen.MoveTo(cow.Aurora(base, line), posx, r.posY(j))
				}
				r.frames <- screen.Flush()
				if k%magic == 0 {
					x++
				}
			}
		} else {
			posx := r.posX(i)

			var n int
			if i > r.screenWidth {
				n = i - r.screenWidth
			}

			base := x * 70
			for j, line := range r.notSaidCowLines {
				y := r.posY(j)
				if i > r.screenWidth {
					if n < len(line) {
						screen.MoveTo(cow.Aurora(base, line[n:]), 1, y)
					} else {
						screen.MoveTo(cow.Aurora(base, " "), 1, y)
					}
				} else if i > len(line) {
					screen.MoveTo(cow.Aurora(base, line), posx, y)
				} else {
					screen.MoveTo(cow.Aurora(base, line[:i]), posx, y)
				}
			}
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
