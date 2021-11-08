package cowsay

import (
	"embed"
	"sort"
	"strings"
)

//go:embed cows/*
var cowsDir embed.FS

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(path string) ([]byte, error) {
	return cowsDir.ReadFile(path)
}

// AssetNames returns the list of filename of the assets.
func AssetNames() []string {
	entries, err := cowsDir.ReadDir("cows")
	if err != nil {
		panic(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".cow")
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

var cowsInBinary = AssetNames()

// CowsInBinary returns the list of cowfiles which are in binary.
// the list is memoized.
func CowsInBinary() []string {
	return cowsInBinary
}
