package cowsay

import (
	"fmt"
	"math"
	"unicode"
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
		if unicode.IsSpace(char) {
			cow.buf.WriteRune(char)
		} else {
			fmt.Fprintf(&cow.buf, "\x1b[%d%sm%c\x1b[0m", rainbow[i%len(rainbow)], attribute, char)
		}
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
		if unicode.IsSpace(char) {
			cow.buf.WriteRune(char)
		} else {
			fmt.Fprintf(&cow.buf, "\033[38;5;%d%sm%c\033[0m", rgb(float64(i)), attribute, char)
		}
		i++
	}
	return cow.buf.String()
}

// https://sking7.github.io/articles/139888127.html#:~:text=value%20of%20frequency.-,Using,-out-of-phase
const (
	freq = 0.01
	m    = math.Pi / 3

	redPhase   = 0
	greenPhase = 2 * m
	bluePhase  = 4 * m
)

var rgbMemo = map[float64]int64{}

func rgb(i float64) int64 {
	if v, ok := rgbMemo[i]; ok {
		return v
	}
	red := int64(6*(math.Sin(freq*i+redPhase)*127+128)/256) * 36
	green := int64(6*(math.Sin(freq*i+greenPhase)*127+128)/256) * 6
	blue := int64(6*(math.Sin(freq*i+bluePhase)*127+128)/256) * 1
	rgbMemo[i] = 16 + red + green + blue
	return rgbMemo[i]
}
