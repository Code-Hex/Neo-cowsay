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
	var attribute string
	if cow.bold {
		attribute = ";1"
	}

	buf := make([]rune, 0, len(mow))
	i := 0
	for _, char := range mow {
		if char == '\n' {
			i = 0
			buf = append(buf, char)
			continue
		}
		buf = append(buf, []rune(fmt.Sprintf("\x1b[%d%sm%c\x1b[0m", rainbow[i%6], attribute, char))...)
		i++
	}

	return string(buf)
}

// Aurora to generate gradation colors string
func (cow *Cow) Aurora(i int, mow string) string {
	var attribute string
	if cow.bold {
		attribute = ";1"
	}

	buf := make([]rune, 0, len(mow))
	for _, char := range mow {
		if char == '\n' {
			buf = append(buf, char)
			continue
		}
		buf = append(buf, []rune(fmt.Sprintf("\033[38;5;%d%sm%c\033[0m", rgb(float64(i)), attribute, char))...)
		i++
	}
	return string(buf)
}

const (
	freq = 0.01
	m    = math.Pi / 3
)

func rgb(i float64) int {
	red := int(6*((math.Sin(freq*i+0)*127+128)/256)) * 36
	green := int(6*((math.Sin(freq*i+2*m)*127+128)/256)) * 6
	blue := int(6*((math.Sin(freq*i+4*m)*127+128)/256)) * 1
	return 16 + red + green + blue
}

// UnANSI color code syntax
func UnANSI(str string) string {
	const maxSearchChars = 30

	l := len(str)
	s := 0
	buf := make([]rune, 0, l)

	for i, char := range str {
		if s > 0 {
			s--
			continue
		}
		if char == '\033' {
			if l > i+1 && str[i+1] == '[' {
				for n := 2; n < maxSearchChars; n++ {
					if l > i+n && str[i+n] == 'm' {
						s = n
						break
					}
				}
			}
			continue
		}

		buf = append(buf, char)
	}

	return string(buf)
}
