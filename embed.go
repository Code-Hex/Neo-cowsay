package cowsay

import (
	"embed"
	"path/filepath"
)

//go:embed cows/*
var cowsDir embed.FS

func Asset(path string) ([]byte, error) {
	return cowsDir.ReadFile(path)
}

func AssetNames() []string {
	const cows = "cows"
	entries, err := cowsDir.ReadDir(cows)
	if err != nil {
		panic(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		filename := filepath.Join(cows, entry.Name())
		names = append(names, filename)
	}
	return names
}
