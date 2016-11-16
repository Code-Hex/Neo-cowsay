package main

import (
	"bufio"
	"fmt"
	"runtime"
	"strings"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	tm "github.com/Code-Hex/goterm"
	colorable "github.com/mattn/go-colorable"
	runewidth "github.com/mattn/go-runewidth"
)

func runSuperCow(cow *cowsay.Cow) error {
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
		return fmt.Errorf("Too height messages...")
	}

	notSaidCow := strings.Split(blank+notSaid, "\n")
	w, cowsWidth := tm.Width(), maxLen(notSaidCow)

	if runtime.GOOS == "windows" {
		tm.Output = bufio.NewWriter(colorable.NewColorableStdout())
	}

	max := w + cowsWidth
	half := max / 2
	x := 0
	for i := 0; i < max+1; i++ {
		tm.Clear()
		if i == half {
			c := make(chan bool)
			go func(t <-chan bool, posx int) {
				for {
					select {
					case <-t:
						return
					default:
						tm.Clear()
						for j, line := range saidCow {
							y := h - len(saidCow) + j - 1
							tm.Println(tm.MoveTo(cow.Aurora(x*70, line), posx, y))
						}
						tm.Flush()
						x++
					}
					time.Sleep(40 * time.Millisecond)
				}
			}(c, w-i)

			time.Sleep(3 * time.Second)
			c <- true
			close(c)
		} else {
			for j, line := range notSaidCow {
				y := h - len(saidCow) + j - 1
				posx := w - i
				if posx < 1 {
					posx = 1
				}
				if i > w {
					n := i - w
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
			time.Sleep(40 * time.Millisecond)
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
