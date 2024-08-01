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

func TestreadFine(t *testing.T) {
	c := readFile("main.go")
	if c == nil {
		t.Error("Context is nil")
	}
	if len(c) < 1 {
		t.Error("return is empty")
	}
}

func TestIsText(t *testing.T) {
	files := [3]string{"main.go", "Makefile", "import-codebase.sh"}
	for _, fname := range files {
		f, _ := os.ReadFile(fname)
		if IsText(f) == false {
			t.Error("mime-type wrong identified")
		}
	}
}
