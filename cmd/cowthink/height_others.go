// +build !windows

package main

import (
	tm "github.com/Code-Hex/goterm"
)

func height() (int, error) {
	return tm.Height(), nil
}
