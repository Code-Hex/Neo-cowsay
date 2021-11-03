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

type line struct {
	text      string
	runeWidth int
}

type lines []*line

func (lines lines) max() int {
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
	}
	return maxWidth
}

func (cow *Cow) getLines(phrase string, width int) lines {
	// Replace tab to 8 spaces
	phrase = strings.Replace(phrase, "\t", "       ", -1)
	text := wordwrap.WrapString(phrase, uint(width))
	lineTexts := strings.Split(text, "\n")
	lines := make([]*line, 0, len(lineTexts))
	for _, lineText := range lineTexts {
		lines = append(lines, &line{
			text:      lineText,
			runeWidth: runewidth.StringWidth(lineText),
		})
	}
	return lines
}

// Balloon to get the balloon and the string entered in the balloon.
func (cow *Cow) Balloon(phrase string) string {
	defer cow.buf.Reset()

	width := cow.ballonWidth
	lines := cow.getLines(phrase, width)
	maxWidth := lines.max()
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
		cow.buf.WriteString(lines[0].text)
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
		cow.padding(lines[i], maxWidth)
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

func (cow *Cow) padding(line *line, maxWidth int) {
	if maxWidth <= line.runeWidth {
		cow.buf.WriteString(line.text)
		return
	}

	cow.buf.WriteString(line.text)
	l := maxWidth - line.runeWidth
	for i := 0; i < l; i++ {
		cow.buf.WriteRune(' ')
	}
}
