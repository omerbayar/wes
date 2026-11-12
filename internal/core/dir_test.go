package core

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListDir(t *testing.T) {
	root := t.TempDir()
	for _, d := range []string{"zeta", "alpha"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"b.txt", "a.txt"} {
		if err := os.WriteFile(filepath.Join(root, f), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := ListDir(root)
	if err != nil {
		t.Fatalf("ListDir: %v", err)
	}

	want := []Entry{
		{Name: "alpha", IsDir: true},
		{Name: "zeta", IsDir: true},
		{Name: "a.txt"},
		{Name: "b.txt"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListDir = %+v, want %+v", got, want)
	}
}

func TestListDirEmpty(t *testing.T) {
	got, err := ListDir(t.TempDir())
	if err != nil {
		t.Fatalf("ListDir: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ListDir = %+v, want empty", got)
	}
}

func TestListDirMissing(t *testing.T) {
	if _, err := ListDir(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("ListDir on missing dir: want error, got nil")
	}
}
