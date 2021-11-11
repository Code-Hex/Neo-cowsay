package screen

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// buffer is the global screen buffer
// Its not recommended write to buffer dirrectly, use package Print,Printf,Println functions instead.
var buffer strings.Builder

// Flush buffer and ensure that it will not overflow screen
func Flush() string {
	defer buffer.Reset()
	return buffer.String()
}

type MoveWriter struct {
	idx  int
	x, y int
	w    io.Writer
	buf  bytes.Buffer
}

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

// MoveTo moves string to position
func MoveTo(str string, x int, y int) {
	x, y = getXY(x, y)
	applyTransform(str, func(idx int, line string) {
		buffer.WriteString("\033[")
		buffer.WriteString(strconv.Itoa(y + idx))
		buffer.WriteRune(';')
		buffer.WriteString(strconv.Itoa(x))
		buffer.WriteRune('H')
		buffer.WriteString(line)
		buffer.WriteString("\033[0K")
	})
}

type sf func(int, string)

// Apply given transformation func for each line in string
func applyTransform(str string, transform sf) {
	for idx, line := range strings.Split(str, "\n") {
		transform(idx, line)
	}
}

// getXY gets relative or absolute coorditantes
// To get relative, set PCT flag to number:
//
//      // Get 10% of total width to `x` and 20 to y
//      x, y = tm.GetXY(10|tm.PCT, 20)
//
func getXY(x int, y int) (int, int) {
	// Set percent flag: num | PCT
	//
	// Check percent flag: num & PCT
	//
	// Reset percent flag: num & 0xFF
	const shift = uint(^uint(0)>>63) << 4
	const PCT = 0x8000 << shift
	if y == -1 {
		y = currentHeight() + 1
	}

	if x&PCT != 0 {
		x = int((x & 0xFF) * Width() / 100)
	}

	if y&PCT != 0 {
		y = int((y & 0xFF) * Height() / 100)
	}

	return x, y
}

// currentHeight returns current height. Line count in Screen buffer.
func currentHeight() int {
	return strings.Count(buffer.String(), "\n")
}
