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

var replacer = strings.NewReplacer("\t", "       ")

func (cow *Cow) getLines(width int) []string {
	// Replace tab to 8 spaces
	cow.phrase = replacer.Replace(cow.phrase)
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

	top, bottom := &strings.Builder{}, &strings.Builder{}
	top.Grow(maxWidth + 2)
	bottom.Grow(maxWidth + 2)

	borderType := cow.borderType()

	top.WriteRune(' ')
	bottom.WriteRune(' ')

	for i := 0; i < maxWidth+2; i++ {
		top.WriteRune('_')
		bottom.WriteRune('_')
	}

	l := len(lines)
	if l == 1 {
		border := borderType.only
		text := fmt.Sprintf("%c %s %c\n", border[0], lines[0], border[1])
		return flush(text, top, bottom)
	}

	var border [2]rune
	var phrase strings.Builder
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		fmt.Fprintf(&phrase, "%c %s %c\n", border[0], padding(lines[i], maxWidth), border[1])
	}

	return flush(phrase.String(), top, bottom)
}

func flush(text string, top, bottom fmt.Stringer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text,
		bottom.String(),
	)
}

func padding(line string, maxWidth int) string {
	w := runewidth.StringWidth(line)
	if maxWidth == w {
		return line
	}

	l := maxWidth - w
	var buf strings.Builder
	buf.Grow(l + len(line))
	buf.WriteString(line)
	for i := 0; i < l; i++ {
		buf.WriteRune(' ')
	}
	return buf.String()
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
