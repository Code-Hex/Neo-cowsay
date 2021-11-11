package decoration

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type options struct {
	withBold    bool
	withRainbow bool
	withAurora  bool
	colorSeq    int
}

func (o *options) maybeBold() string {
	if o.withBold {
		return ";1"
	}
	return ""
}

// Option for any writer in this package.
type Option func(o *options)

// WithBold writes with bold.
func WithBold() Option {
	return func(o *options) {
		o.withBold = true
	}
}

// WithRainbow writes with rainbow.
func WithRainbow() Option {
	return func(o *options) {
		o.withRainbow = true
	}
}

// WithAurora writes with aurora.
func WithAurora(initialSeq int) Option {
	return func(o *options) {
		o.withAurora = true
		o.colorSeq = initialSeq
	}
}

// Writer is a writer to decorates.
type Writer struct {
	writer  io.Writer
	buf     bytes.Buffer
	options *options
}

// NewWriter creates a new writer.
func NewWriter(w io.Writer, opts ...Option) *Writer {
	options := new(options)
	for _, optFunc := range opts {
		optFunc(options)
	}
	return &Writer{
		writer:  w,
		options: options,
	}
}

// SetColorSeq sets current color sequence.
func (w *Writer) SetColorSeq(colorSeq int) {
	w.options.colorSeq = colorSeq
}

func (w *Writer) Write(b []byte) (nn int, err error) {
	switch {
	case w.options.withAurora:
		return w.writeAsAurora(b)
	case w.options.withRainbow:
		return w.writeAsRainbow(b)
	case w.options.withBold:
		return w.writeAsDefaultBold(b)
	default:
		return w.writer.Write(b)
	}
}

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
