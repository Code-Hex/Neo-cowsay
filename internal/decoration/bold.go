package decoration

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

func (w *Writer) writeAsDefaultBold(b []byte) (nn int, err error) {
	defer w.buf.Reset()

	for len(b) > 0 {
		char, size := utf8.DecodeRune(b)
		if char == '\n' {
			w.buf.WriteRune(char)
			b = b[size:]
			continue
		}
		if unicode.IsSpace(char) {
			w.buf.WriteRune(char)
		} else {
			fmt.Fprintf(&w.buf, "\x1b[1m%c\x1b[0m", char)
		}
		b = b[size:]
	}
	return w.writer.Write(w.buf.Bytes())
}

func (w *Writer) writeStringAsDefaultBold(s string) (nn int, err error) {
	defer w.buf.Reset()

	for _, char := range s {
		if char == '\n' {
			w.buf.WriteRune(char)
			continue
		}
		if unicode.IsSpace(char) {
			w.buf.WriteRune(char)
		} else {
			fmt.Fprintf(&w.buf, "\x1b[1m%c\x1b[0m", char)
		}
	}
	if sw, ok := w.writer.(io.StringWriter); ok {
		return sw.WriteString(w.buf.String())
	}
	return w.writer.Write(w.buf.Bytes())
}
