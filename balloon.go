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

func (cow *Cow) getLines(width int) []string {
	// Replace tab to 8 spaces
	cow.phrase = strings.Replace(cow.phrase, "\t", "       ", -1)
	text := wordwrap.WrapString(cow.phrase, uint(width))
	return strings.Split(text, "\n")
}

// Balloon to get the balloon and the string entered in the balloon.
func (cow *Cow) Balloon() string {
	defer cow.buf.Reset()

	width := cow.ballonWidth
	lines := cow.getLines(width)
	// find max length from text lines
	maxWidth := max(lines)
	if maxWidth > width {
		maxWidth = width
	}

	top := make([]byte, 0, maxWidth+2)
	bottom := make([]byte, 0, maxWidth+2)

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
		cow.buf.Write(top)
		cow.buf.WriteRune('\n')
		cow.buf.WriteRune(border[0])
		cow.buf.WriteRune(' ')
		cow.buf.WriteString(lines[0])
		cow.buf.WriteRune(' ')
		cow.buf.WriteRune(border[1])
		cow.buf.WriteRune('\n')
		cow.buf.Write(bottom)
		cow.buf.WriteRune('\n')
		return cow.buf.String()
	}

	cow.buf.Write(top)
	cow.buf.WriteRune('\n')

	var border [2]rune
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		cow.buf.WriteRune(border[0])
		cow.buf.WriteRune(' ')
		padding(&cow.buf, lines[i], maxWidth)
		cow.buf.WriteRune(' ')
		cow.buf.WriteRune(border[1])
		cow.buf.WriteRune('\n')
	}

	cow.buf.Write(bottom)
	cow.buf.WriteRune('\n')
	return cow.buf.String()
}

func (cow *Cow) flush(text, top, bottom fmt.Stringer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text.String(),
		bottom.String(),
	)
}

func padding(b *strings.Builder, line string, maxWidth int) {
	w := runewidth.StringWidth(line)
	if maxWidth <= w {
		b.WriteString(line)
		return
	}

	l := maxWidth - w
	for i := 0; i < l; i++ {
		b.WriteRune(' ')
	}
	b.WriteString(line)
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
