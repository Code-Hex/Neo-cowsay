package cowsay

import (
	"bytes"
	"fmt"
	"math"
)

const (
	red = iota + 31
	green
	yellow
	blue
	magenta
	cyan
)

// Rainbow to generate rainbow string
func (cow *Cow) Rainbow(mow string) string {
	var attribute string
	if cow.Bold {
		attribute = ";1"
	}

	rainbow := []int{magenta, red, yellow, green, cyan, blue}
	b := bytes.NewBuffer(make([]byte, 0, len(mow)))
	i := 0
	for _, char := range mow {
		if char == '\n' {
			i = 0
			b.WriteRune(char)
			continue
		}
		b.WriteString(fmt.Sprintf("\x1b[%d%sm%c\x1b[0m", rainbow[i%6], attribute, char))
		i++
	}

	return b.String()
}

// Aurora to generate gradation colors string
func (cow *Cow) Aurora(i int, mow string) string {
	var attribute string
	if cow.Bold {
		attribute = ";1"
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(mow)))
	for _, char := range mow {
		if char == '\n' {
			buf.WriteRune(char)
			continue
		}

		buf.WriteString(fmt.Sprintf("\033[38;5;%d%sm%c\033[0m", rgb(float64(i)), attribute, char))
		i++
	}
	return buf.String()
}

func rgb(i float64) int {
	freq := 0.01
	red := int(6*((math.Sin(freq*i+0)*127+128)/256)) * 36
	green := int(6*((math.Sin(freq*i+2*math.Pi/3)*127+128)/256)) * 6
	blue := int(6*((math.Sin(freq*i+4*math.Pi/3)*127+128)/256)) * 1

	return 16 + red + green + blue
}
