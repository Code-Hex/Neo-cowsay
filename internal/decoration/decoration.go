package decoration

import (
	"bytes"
	"io"
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

var _ interface {
	io.Writer
	io.StringWriter
} = (*Writer)(nil)

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

// Write writes bytes. which is implemented io.Writer.
//
// If Bold is enabled in the options, the text will be written as Bold.
// If both Aurora and Rainbow are enabled in the options, Aurora will take precedence.
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

// WriteString writes string. which is implemented io.StringWriter.
//
// See also Write.
func (w *Writer) WriteString(s string) (n int, err error) {
	switch {
	case w.options.withAurora:
		return w.writeStringAsAurora(s)
	case w.options.withRainbow:
		return w.writeStringAsRainbow(s)
	case w.options.withBold:
		return w.writeStringAsDefaultBold(s)
	default:
		if sw, ok := w.writer.(io.StringWriter); ok {
			return sw.WriteString(w.buf.String())
		}
		return w.writer.Write([]byte(s))
	}
}
