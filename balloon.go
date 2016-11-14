package cowsay

import (
	"bytes"
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
	if cow.Thinking {
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
	r := strings.NewReplacer("\t", "       ")
	cow.Phrase = r.Replace(cow.Phrase)
	text := wordwrap.WrapString(cow.Phrase, uint(width))
	return strings.Split(text, "\n")
}

func (cow *Cow) balloon() string {
	width := cow.BallonWidth
	if width <= 0 {
		width = 1
		cow.Phrase = "0"
	}

	lines := cow.getLines(width)
	// find max length from text lines
	maxWidth := max(lines)
	if maxWidth > width {
		maxWidth = width
	}

	top := bytes.NewBuffer(make([]byte, 0, maxWidth+2))
	bottom := bytes.NewBuffer(make([]byte, 0, maxWidth+2))

	borderType := cow.borderType()

	top.WriteRune(' ')
	bottom.WriteRune(' ')

	for i := 0; i < maxWidth+2; i++ {
		top.WriteRune('_')
		bottom.WriteRune('-')
	}

	l := len(lines)
	if l == 1 {
		border := borderType.only
		return flush(fmt.Sprintf("%c %s %c\n", border[0], lines[0], border[1]), top, bottom)
	}

	var border [2]rune
	var phrase bytes.Buffer
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		phrase.WriteString(fmt.Sprintf("%c %s %c\n", border[0], padding(lines[i], maxWidth), border[1]))
	}

	return flush(phrase.String(), top, bottom)
}

func flush(text string, top, bottom *bytes.Buffer) string {
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
	pad := make([]rune, 0, l)
	for i := 0; i < l; i++ {
		pad = append(pad, ' ')
	}
	return line + string(pad)
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
