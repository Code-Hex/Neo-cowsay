package super

import (
	"errors"
	"io"
	"os"
	"strings"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	tm "github.com/Code-Hex/Neo-cowsay/internal/screen"
	colorable "github.com/mattn/go-colorable"
	runewidth "github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	span    = 20 * time.Millisecond
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
	w, cowsWidth := tm.Width(), maxLen(notSaidCow)
	output := colorable.NewColorableStdout()

	output.Write([]byte("\033[2J\033[?25l"))

	// tm.Clear()

	views := make(chan string, 10000)
	go func() {
		max := w + cowsWidth
		half := max / 2
		diff := h - len(saidCow)
		for x, i := 0, 0; i <= max; i++ {
			if i == half {
				posx := w - i
				times := standup / span
				for k := 0; k < int(times); k++ {
					// draw colored cow
					base := x * 70
					for j, line := range saidCow {
						y := diff + j - 1
						tm.MoveTo(cow.Aurora(base, line), posx, y)
					}
					views <- tm.Flush()
					x++
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
							tm.MoveTo(cow.Aurora(base, line[n:]), 1, y)
						} else {
							tm.MoveTo(cow.Aurora(base, " "), 1, y)
						}
					} else if i > len(line) {
						tm.MoveTo(cow.Aurora(base, line), posx, y)
					} else {
						tm.MoveTo(cow.Aurora(base, line[:i]), posx, y)
					}
				}
				views <- tm.Flush()
			}
			x++
		}
		close(views)
	}()

	for view := range views {
		io.Copy(output, strings.NewReader(view))
		time.Sleep(span)
	}

	output.Write([]byte("\033[?25h"))

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
