package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	wordwrap "github.com/mitchellh/go-wordwrap"
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

func (cow *Cow) balloon(width int) string {
	text := wordwrap.WrapString(cow.Phrase, uint(width))
	lines := strings.Split(text, "\n")

	// find max length from text lines
	maxWidth := max(lines)
	if maxWidth > width {
		maxWidth = width
	}

	var (
		top    bytes.Buffer
		bottom bytes.Buffer
	)

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

func flush(text string, top, bottom bytes.Buffer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text,
		bottom.String(),
	)
}

func padding(line string, maxWidth int) string {
	return line + strings.Repeat(" ", maxWidth-len(line))
}

func max(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		len := utf8.RuneCountInString(line)
		if len > maxWidth {
			maxWidth = len
		}
	}
	return maxWidth
}
