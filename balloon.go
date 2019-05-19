package cowsay

import (
	"fmt"
	"strings"

	wordwrap "github.com/Code-Hex/go-wordwrap"
	runewidth "github.com/mattn/go-runewidth"
)

type border struct {
	first  [2]rune
	middle [2]rune
	last   [2]rune
	only   [2]rune
}

type buf []rune

func (b buf) String() string {
	return string(b)
}

func (cow *Cow) borderType() border {
	if cow.thinking {
		return border{
			first:  [2]rune{'(', ')'},
			middle: [2]rune{'(', ')'},
			last:   [2]rune{'(', ')'},
			only:   [2]rune{'(', ')'},
		}
	}

	return border{
		first:  [2]rune{'/', '\\'},
		middle: [2]rune{'|', '|'},
		last:   [2]rune{'\\', '/'},
		only:   [2]rune{'<', '>'},
	}
}

func (cow *Cow) getLines(width int) []string {
	// Replace tab to 8 spaces
	cow.phrase = strings.Replace(cow.phrase, "\t", "       ", -1)
	text := wordwrap.WrapString(cow.phrase, uint(width))
	return strings.Split(text, "\n")
}

// Balloon to get the balloon and the string entered in the balloon.
func (cow *Cow) Balloon() string {
	width := cow.ballonWidth
	lines := cow.getLines(width)
	// find max length from text lines
	maxWidth := max(lines)
	if maxWidth > width {
		maxWidth = width
	}

	top := make(buf, 0, maxWidth+2)
	bottom := make(buf, 0, maxWidth+2)

	borderType := cow.borderType()

	top = append(top, ' ')
	bottom = append(bottom, ' ')

	for i := 0; i < maxWidth+2; i++ {
		top = append(top, '_')
		bottom = append(bottom, '-')
	}

	l := len(lines)
	if l == 1 {
		border := borderType.only
		text := fmt.Sprintf("%c %s %c\n", border[0], lines[0], border[1])
		return flush(buf(text), top, bottom)
	}

	var border [2]rune
	phrase := make(buf, 0, l*100)
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		phrase = append(phrase, buf(fmt.Sprintf("%c %s %c\n", border[0], padding(lines[i], maxWidth), border[1]))...)
	}
	return flush(phrase, top, bottom)
}

func flush(text, top, bottom fmt.Stringer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text.String(),
		bottom.String(),
	)
}

func padding(line string, maxWidth int) string {
	w := runewidth.StringWidth(line)
	if maxWidth == w {
		return line
	}

	l := maxWidth - w
	pad := make(buf, l)
	for i := 0; i < l; i++ {
		pad[i] = ' '
	}
	return line + pad.String()
}

func max(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		len := runewidth.StringWidth(line)
		if len > maxWidth {
			maxWidth = len
		}
	}
	return maxWidth
}
