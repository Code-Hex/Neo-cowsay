package decoration

import (
	"fmt"
	"io"
	"math"
	"unicode"
	"unicode/utf8"
)

func (w *Writer) writeAsAurora(b []byte) (nn int, err error) {
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
			_, err := fmt.Fprintf(&w.buf, "\033[38;5;%d%sm%c\033[0m",
				rgb(float64(w.options.colorSeq)),
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

func (w *Writer) writeStringAsAurora(s string) (nn int, err error) {
	defer w.buf.Reset()

	for _, char := range s {
		if char == '\n' {
			w.buf.WriteRune(char)
			continue
		}
		if unicode.IsSpace(char) {
			w.buf.WriteRune(char)
		} else {
			_, err := fmt.Fprintf(&w.buf, "\033[38;5;%d%sm%c\033[0m",
				rgb(float64(w.options.colorSeq)),
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

// https://sking7.github.io/articles/139888127.html#:~:text=value%20of%20frequency.-,Using,-out-of-phase
const (
	freq = 0.01
	m    = math.Pi / 3

	redPhase   = 0
	greenPhase = 2 * m
	bluePhase  = 4 * m
)

var rgbMemo = map[float64]int64{}

func rgb(i float64) int64 {
	if v, ok := rgbMemo[i]; ok {
		return v
	}
	red := int64(6*(math.Sin(freq*i+redPhase)*127+128)/256) * 36
	green := int64(6*(math.Sin(freq*i+greenPhase)*127+128)/256) * 6
	blue := int64(6*(math.Sin(freq*i+bluePhase)*127+128)/256) * 1
	rgbMemo[i] = 16 + red + green + blue
	return rgbMemo[i]
}
