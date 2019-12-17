package super

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	"github.com/Code-Hex/Neo-cowsay/internal/screen"
	runewidth "github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
)

// Frequency the color changes
const magic = 2

const (
	span    = 30 * time.Millisecond
	standup = 3 * time.Second
)

// RunSuperCow runs super cow mode animation on the your terminal
func RunSuperCow(cow *cowsay.Cow) error {
	balloon := cow.Balloon()
	blank := createBlankSpace(balloon)

	said, err := cow.GetCow(0)
	if err != nil {
		return err
	}

	notSaid, err := cow.GetCow(' ')
	if err != nil {
		return err
	}

	saidCow := strings.Split(balloon+said, "\n")

	// When it is higher than the height of the terminal
	h, err := height()
	if err != nil {
		return err
	}
	if len(saidCow) > h {
		return errors.New("Too height messages")
	}

	notSaidCow := strings.Split(blank+notSaid, "\n")
	w, cowsWidth := screen.Width(), maxLen(notSaidCow)

	max := w + cowsWidth
	half := max / 2
	diff := h - len(saidCow)
	views := make(chan string, max)

	screen.SaveState()
	screen.HideCursor()
	screen.Clear()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for x, i := 0, 0; i <= max; i++ {
			if i == half {
				posx := w - i
				times := standup / span
				for k := 0; k < int(times); k++ {
					// draw colored cow
					base := x * 70
					for j, line := range saidCow {
						y := diff + j - 1
						screen.MoveTo(cow.Aurora(base, line), posx, y)
					}
					views <- screen.Flush()
					if k%magic == 0 {
						x++
					}
				}
			} else {
				posx := w - i
				if posx < 1 {
					posx = 1
				}

				var n int
				if i > w {
					n = i - w
				}

				base := x * 70
				for j, line := range notSaidCow {
					y := diff + j - 1
					if i > w {
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
				views <- screen.Flush()
			}
			if i%magic == 0 {
				x++
			}
		}
		close(views)
	}()

LOOP:
	for view := range views {
		select {
		case <-quit:
			screen.Clear()
			break LOOP
		default:
		}
		screen.Stdout.Write([]byte(view))
		time.Sleep(span)
	}

	screen.UnHideCursor()
	screen.RestoreState()

	return nil
}

var buf strings.Builder

func createBlankSpace(balloon string) string {
	buf.Reset()
	l := len(strings.Split(balloon, "\n"))
	for i := 1; i < l; i++ {
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

func height() (int, error) {
	fd := int(os.Stdout.Fd())
	_, height, err := terminal.GetSize(fd)
	if err != nil {
		return 0, err
	}
	return height, nil
}
