//go:build !windows
// +build !windows

package cowsay

import "strings"

func splitPath(s string) []string {
	return strings.Split(s, ":")
}
