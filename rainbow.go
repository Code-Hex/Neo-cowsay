package main

import (
	"bytes"
	"fmt"
)

const (
	red = iota + 31
	green
	yellow
	blue
	magenta
	cyan
)

func makeRainbow(mow string) string {
	rainbow := []int{magenta, red, yellow, green, cyan, blue}
	b := bytes.NewBuffer(make([]byte, 0, len(mow)))
	i := 0
	for _, char := range mow {
		if char == '\n' {
			i = 0
			b.WriteRune('\n')
			continue
		}
		b.WriteString(fmt.Sprintf("\x1b[%dm%c\x1b[0m", rainbow[i%6], char))
		i++
	}

	return b.String()
}
