package screen

import (
	"os"
	"sync"

	"github.com/mattn/go-colorable"
	"golang.org/x/crypto/ssh/terminal"
)

// Stdout color supported stdout
var Stdout = colorable.NewColorableStdout()

// SaveState saves cursor state.
func SaveState() { Stdout.Write([]byte("\0337")) }

// RestoreState restores cursor state.
func RestoreState() { Stdout.Write([]byte("\0338")) }

// Clear clears terminal screen.
func Clear() { Stdout.Write([]byte("\033[2J")) }

// HideCursor hide the cursor
func HideCursor() { Stdout.Write([]byte("\033[?25l")) }

// UnHideCursor unhide the cursor
func UnHideCursor() { Stdout.Write([]byte("\033[?25h")) }

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
