package screen

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
)

// Output means buffered stdout
var Output *bufio.Writer = NewWriter(os.Stdout)

// Screen is the global screen buffer
// Its not recommented write to buffer dirrectly, use package Print,Printf,Println fucntions instead.
var Screen *strings.Builder = new(strings.Builder)

// NewWriter creates buffered writer
func NewWriter(w io.Writer) *bufio.Writer {
	width, height := Width(), Height()
	if width != -1 && height != -1 {
		return bufio.NewWriterSize(w, width*height+20000)
	}
	return bufio.NewWriter(w)
}

// Clear screen
func Clear() {
	Output.WriteString("\033[2J")
}

// Flush buffer and ensure that it will not overflow screen
func Flush() string {
	defer Screen.Reset()
	return Screen.String()
}

var size struct {
	once   sync.Once
	width  int
	height int
}

func getSize() (int, int) {
	size.once.Do(func() {
		var err error
		size.width, size.height, err = terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			size.width, size.height = -1, -1
		}
	})
	return size.width, size.height
}

// Width returns console width
func Width() int {
	width, _ := getSize()
	return width
}

// Height returns console height
func Height() int {
	_, height := getSize()
	return height
}

// MoveTo moves string to possition
func MoveTo(str string, x int, y int) {
	x, y = getXY(x, y)
	applyTransform(str, func(idx int, line string) {
		Screen.WriteString("\033[")
		Screen.WriteString(strconv.Itoa(y + idx))
		Screen.WriteRune(';')
		Screen.WriteString(strconv.Itoa(x))
		Screen.WriteRune('H')
		Screen.WriteString(line)
		Screen.WriteString("\033[0K")
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
	return strings.Count(Screen.String(), "\n")
}
