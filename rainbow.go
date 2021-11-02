package cowsay

import (
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

var rainbow = []int{magenta, red, yellow, green, cyan, blue}

// Rainbow to generate rainbow string
func (cow *Cow) Rainbow(mow string) string {
	defer cow.buf.Reset()

	var attribute string
	if cow.bold {
		attribute = ";1"
	}

	i := 0
	for _, char := range mow {
		if char == '\n' {
			i = 0
			cow.buf.WriteRune(char)
			continue
		}
		fmt.Fprintf(&cow.buf, "\x1b[%d%sm%c\x1b[0m", rainbow[i%6], attribute, char)
		i++
	}

	return cow.buf.String()
}

// Aurora to generate gradation colors string
func (cow *Cow) Aurora(i int, mow string) string {
	defer cow.buf.Reset()

	var attribute string
	if cow.bold {
		attribute = ";1"
	}

	for _, char := range mow {
		if char == '\n' {
			cow.buf.WriteRune(char)
			continue
		}
		fmt.Fprintf(&cow.buf, "\033[38;5;%d%sm%c\033[0m", rgb(float64(i)), attribute, char)
		i++
	}
	return cow.buf.String()
}

const (
	freq = 0.01
	m    = math.Pi / 3
)

func rgb(i float64) int64 {
	red := int64(6*((math.Sin(freq*i+0)*127+128)/256)) * 36
	green := int64(6*((math.Sin(freq*i+2*m)*127+128)/256)) * 6
	blue := int64(6*((math.Sin(freq*i+4*m)*127+128)/256)) * 1
	return 16 + red + green + blue
}
