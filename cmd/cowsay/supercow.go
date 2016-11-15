package main

import (
	"fmt"
	"strings"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	tm "github.com/Code-Hex/goterm"
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
	h := tm.Height()
	if len(saidCow) > h {
		return fmt.Errorf("Too height messages...")
	}

	notSaidCow := strings.Split(blank+notSaid, "\n")
	w, cowsWidth := tm.Width(), maxLen(notSaidCow)

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
						for _, line := range saidCow {
							tm.Println(tm.MoveTo(cow.Aurora(x*70, line), posx, h))
						}
						tm.Flush()
						x++
					}
					time.Sleep(40 * time.Millisecond)
				}
			}(c, w-i)

			time.Sleep(3 * time.Second)
			c <- true
		} else {
			for _, line := range notSaidCow {
				posx := w - i
				if i > w {
					n := i - w
					if n < len(line) {
						tm.Println(tm.MoveTo(cow.Aurora(x*70, line[n:]), 0, h))
					} else {
						tm.Println(" ")
					}
				} else if i > len(line) {
					tm.Println(tm.MoveTo(cow.Aurora(x*70, line), posx, h))
				} else {
					tm.Println(tm.MoveTo(cow.Aurora(x*70, line[:i]), posx, h))
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
