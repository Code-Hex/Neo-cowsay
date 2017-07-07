package super

import (
	"bufio"
	"errors"
	"runtime"
	"strings"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	tm "github.com/Code-Hex/goterm"
	colorable "github.com/mattn/go-colorable"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	span    = 40 * time.Millisecond
	standup = 3 * time.Second
)

func RunSuperCow(cow *cowsay.Cow) error {
	cowsay.CowsInit(cow)
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

	if runtime.GOOS == "windows" {
		tm.Output = bufio.NewWriter(colorable.NewColorableStdout())
	}

	max := w + cowsWidth
	half := max / 2
	diff := h - len(saidCow)
	for x, i := 0, 0; i <= max; i++ {
		tm.Clear()
		if i == half {
			posx := w - i
			after := time.After(standup)
		DRAW:
			for {
				select {
				case <-after:
					break DRAW
				default:
					tm.Clear()
					// draw colored cow
					for j, line := range saidCow {
						y := diff + j - 1
						tm.Println(tm.MoveTo(cow.Aurora(x*70, line), posx, y))
					}
					tm.Flush()
					x++
				}
				time.Sleep(span)
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

			for j, line := range notSaidCow {
				y := diff + j - 1
				if i > w {
					if n < len(line) {
						tm.Print(tm.MoveTo(cow.Aurora(x*70, line[n:]), 1, y))
					} else {
						tm.Print(tm.MoveTo(cow.Aurora(x*70, " "), 1, y))
					}
				} else if i > len(line) {
					tm.Print(tm.MoveTo(cow.Aurora(x*70, line), posx, y))
				} else {
					tm.Print(tm.MoveTo(cow.Aurora(x*70, line[:i]), posx, y))
				}
			}
			tm.Flush()
			time.Sleep(span)
		}
		x++
	}

	return nil
}

func createBlankSpace(balloon string) string {
	l := len(strings.Split(balloon, "\n"))
	blank := make([]byte, 0, l)
	for i := 1; i < l; i++ {
		blank = append(blank, byte('\n'))
	}
	return string(blank)
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
