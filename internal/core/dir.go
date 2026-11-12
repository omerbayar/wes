// Package core holds wes's UI-independent logic.
package core

import (
	"os"
	"sort"
)

// Entry is a single item inside a directory.
type Entry struct {
	Name  string
	IsDir bool
}

// ListDir returns the entries of the directory at path, directories first,
// each group sorted by name.
func ListDir(path string) ([]Entry, error) {
	des, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(des))
	for _, de := range des {
		entries = append(entries, Entry{Name: de.Name(), IsDir: de.IsDir()})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}
