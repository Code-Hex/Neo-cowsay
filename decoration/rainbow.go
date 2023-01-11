package decoration

import (
	"fmt"
	"io"
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
			_, err := fmt.Fprintf(&w.buf, "\x1b[%d%sm%c\x1b[0m",
				rainbow[w.options.colorSeq%len(rainbow)],
				w.options.maybeBold(),
				char,
			)
			if err != nil {
				return 0, err
			}
		}
		w.options.colorSeq++
		b = b[size:]
	}

	return w.writer.Write(w.buf.Bytes())
}

func (w *Writer) writeStringAsRainbow(s string) (nn int, err error) {
	defer w.buf.Reset()

	for _, char := range s {
		if char == '\n' {
			w.options.colorSeq = 0
			w.buf.WriteRune(char)
			continue
		}
		if unicode.IsSpace(char) {
			w.buf.WriteRune(char)
		} else {
			_, err := fmt.Fprintf(&w.buf, "\x1b[%d%sm%c\x1b[0m",
				rainbow[w.options.colorSeq%len(rainbow)],
				w.options.maybeBold(),
				char,
			)
			if err != nil {
				return 0, err
			}
		}
		w.options.colorSeq++
	}

	if sw, ok := w.writer.(io.StringWriter); ok {
		return sw.WriteString(w.buf.String())
	}
	return w.writer.Write(w.buf.Bytes())
}
