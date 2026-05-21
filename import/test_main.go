package main

import (
	"os"
	"testing"
)

func TestIsFile(t *testing.T) {
	if IsFile("main.go") != true {
		t.Error("IsFile should return true")
	}
	if IsFile("foobar.go") == true {
		t.Error("IsFile should return false")
	}
}

func TestReadFile(t *testing.T) {
	c, err := readFile("main.go")
	if err != nil {
		t.Fatalf("readFile error: %v", err)
	}
	if c == nil {
		t.Error("Context is nil")
	}
	if len(c) < 1 {
		t.Error("return is empty")
	}
}

func TestIsText(t *testing.T) {
	files := [3]string{"main.go", "Makefile", "import/main.go"}
	for _, fname := range files {
		f, err := os.ReadFile(fname)
		if err != nil {
			t.Fatalf("readFile(%s) error: %v", fname, err)
		}
		if IsText(f) == false {
			t.Errorf("mime-type wrong identified for %s", fname)
		}
	}
}
