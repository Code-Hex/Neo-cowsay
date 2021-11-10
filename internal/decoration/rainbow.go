package decoration

import (
	"fmt"
	"unicode"
	"unicode/utf8"
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

func (w *Writer) writeAsRainbow(b []byte) (nn int, err error) {
	defer w.buf.Reset()

	for len(b) > 0 {
		char, size := utf8.DecodeRune(b)
		if char == '\n' {
			w.options.colorSeq = 0
			w.buf.WriteRune(char)
			b = b[size:]
			continue
		}
		if unicode.IsSpace(char) {
			w.buf.WriteRune(char)
		} else {
			fmt.Fprintf(&w.buf, "\x1b[%d%sm%c\x1b[0m",
				rainbow[w.options.colorSeq%len(rainbow)],
				w.options.maybeBold(),
				char,
			)
		}
		w.options.colorSeq++
		b = b[size:]
	}

	return w.writer.Write(w.buf.Bytes())
}
