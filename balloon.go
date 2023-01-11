package cowsay

import (
	"fmt"
	"strings"

	"github.com/Code-Hex/go-wordwrap"
	"github.com/mattn/go-runewidth"
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

func (cow *Cow) maxLineWidth(lines []*line) int {
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
		if !cow.disableWordWrap && maxWidth > cow.ballonWidth {
			return cow.ballonWidth
		}
	}
	return maxWidth
}

func (cow *Cow) getLines(phrase string) []*line {
	text := cow.canonicalizePhrase(phrase)
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

func (cow *Cow) canonicalizePhrase(phrase string) string {
	// Replace tab to 8 spaces
	phrase = strings.Replace(phrase, "\t", "       ", -1)

	if cow.disableWordWrap {
		return phrase
	}
	width := cow.ballonWidth
	return wordwrap.WrapString(phrase, uint(width))
}

// Balloon to get the balloon and the string entered in the balloon.
func (cow *Cow) Balloon(phrase string) string {
	defer cow.buf.Reset()

	lines := cow.getLines(phrase)
	maxWidth := cow.maxLineWidth(lines)

	cow.writeBallon(lines, maxWidth)

	return cow.buf.String()
}

func (cow *Cow) writeBallon(lines []*line, maxWidth int) {
	top := make([]byte, 0, maxWidth+2)
	bottom := make([]byte, 0, maxWidth+2)

	top = append(top, ' ')
	bottom = append(bottom, ' ')

	for i := 0; i < maxWidth+2; i++ {
		top = append(top, '_')
		bottom = append(bottom, '-')
	}

	borderType := cow.borderType()

	cow.buf.Write(top)
	cow.buf.Write([]byte{' ', '\n'})
	defer func() {
		cow.buf.Write(bottom)
		cow.buf.Write([]byte{' ', '\n'})
	}()

	l := len(lines)
	if l == 1 {
		border := borderType.only
		cow.buf.WriteRune(border[0])
		cow.buf.WriteRune(' ')
		cow.buf.WriteString(lines[0].text)
		cow.buf.WriteRune(' ')
		cow.buf.WriteRune(border[1])
		cow.buf.WriteRune('\n')
		return
	}

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
