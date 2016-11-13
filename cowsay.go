package cow

import (
	"bytes"
	"math/rand"
	"strings"
	"unicode/utf8"

	wordwrap "github.com/mitchellh/go-wordwrap"
)

type Cow struct {
	Phrase   string
	Eyes     [2]rune
	Tongue   [2]rune
	Random   bool
	Type     string
	Thinking bool
}

func Say(cow *Cow) (string, error) {

	if cow.Random {
		cow.Type = PickCow()
	}

	if cow.Type == "" {
		cow.Type = "cows/default.cow"
	}

	if !strings.HasSuffix(cow.Type, ".cow") {
		cow.Type += ".cow"
	}

	if !strings.HasPrefix(cow.Type, "cows/") {
		cow.Type = "cows/" + cow.Type
	}

}

func PickCow() string {
	var cows []string
	for _, path := range _bindata {
		cows = append(cows, path)
	}

	return pickup(cows)
}

func pickup(cows []string) string {
	n := len(cows)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		cows[i], cows[j] = cows[j], cows[i]
	}
	return cows[rand.Intn(n)]
}

type border struct {
	first  [2]rune
	middle [2]rune
	last   [2]rune
	only   [2]rune
}

func (cow *Cow) constructBalloon(width int) {
	text = wordwrap.WrapString(cow.Phrase, uint(width))
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

	for i := 0; i < maxWidth; i++ {
		top.WriteRune('_')
		bottom.WriteRune('-')
	}

	l := len(line)
	border := make([]rune, 2)
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
	}

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
}

func (cow *Cow) borderType() border {
	if cow.Thinking {
		return border{
			first:  {'(', ')'},
			middle: {'(', ')'},
			last:   {'(', ')'},
			only:   {'(', ')'},
		}
	}

	return border{
		first: {'/', '\\'},
		midle: {'|', '|'},
		last:  {'\\', '/'},
		only:  {'<', '>'},
	}
}
