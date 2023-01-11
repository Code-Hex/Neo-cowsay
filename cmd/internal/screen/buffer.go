package screen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// buffer is the global screen buffer
// It's not recommended write to buffer directly, use package Print,Printf,Println functions instead.
var buffer strings.Builder

// Flush buffer and ensure that it will not overflow screen
func Flush() string {
	defer buffer.Reset()
	return buffer.String()
}

// MoveWriter is implemented io.Writer and io.StringWriter.
type MoveWriter struct {
	idx  int
	x, y int
	w    io.Writer
	buf  bytes.Buffer
}

var _ interface {
	io.Writer
	io.StringWriter
} = (*MoveWriter)(nil)

// NewMoveWriter creates a new MoveWriter.
func NewMoveWriter(w io.Writer, x, y int) *MoveWriter {
	x, y = getXY(x, y)
	return &MoveWriter{
		w: w,
		x: x,
		y: y,
	}
}

// SetPosx sets pos x
func (m *MoveWriter) SetPosx(x int) {
	x, _ = getXY(x, 0)
	m.x = x
}

// Reset resets
func (m *MoveWriter) Reset() {
	m.idx = 0
	m.buf.Reset()
}

// Write writes bytes. which is implemented io.Writer.
func (m *MoveWriter) Write(bs []byte) (nn int, _ error) {
	br := bytes.NewReader(bs)
	for {
		b, err := br.ReadByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			n, _ := fmt.Fprintf(m.w, "\x1b[%d;%dH%s\x1b[0K",
				m.y+m.idx,
				m.x,
				m.buf.String(),
			)
			nn += n
			return
		}
		if b == '\n' {
			n, _ := fmt.Fprintf(m.w, "\x1b[%d;%dH%s\x1b[0K",
				m.y+m.idx,
				m.x,
				m.buf.String(),
			)
			m.buf.Reset()
			m.idx++
			nn += n
		} else {
			m.buf.WriteByte(b)
		}
	}
}

// WriteString writes string. which is implemented io.StringWriter.
func (m *MoveWriter) WriteString(s string) (nn int, _ error) {
	for _, char := range s {
		if char == '\n' {
			n, _ := fmt.Fprintf(m.w, "\x1b[%d;%dH%s\x1b[0K",
				m.y+m.idx,
				m.x,
				m.buf.String(),
			)
			m.buf.Reset()
			m.idx++
			nn += n
		} else {
			m.buf.WriteRune(char)
		}
	}
	if m.buf.Len() > 0 {
		n, _ := fmt.Fprintf(m.w, "\x1b[%d;%dH%s\x1b[0K",
			m.y+m.idx,
			m.x,
			m.buf.String(),
		)
		nn += n
	}
	return
}

// getXY gets relative or absolute coorditantes
// To get relative, set PCT flag to number:
//
//	// Get 10% of total width to `x` and 20 to y
//	x, y = tm.GetXY(10|tm.PCT, 20)
func getXY(x int, y int) (int, int) {
	// Set percent flag: num | PCT
	//
	// Check percent flag: num & PCT
	//
	// Reset percent flag: num & 0xFF
	const shift = ^uint(0) >> 63 << 4
	const PCT = 0x8000 << shift
	if y == -1 {
		y = currentHeight() + 1
	}

	if x&PCT != 0 {
		x = (x & 0xFF) * Width() / 100
	}

	if y&PCT != 0 {
		y = (y & 0xFF) * Height() / 100
	}

	return x, y
}

// currentHeight returns current height. Line count in Screen buffer.
func currentHeight() int {
	return strings.Count(buffer.String(), "\n")
}
