package main

import (
	"os"
	"testing"
)

func TestProgram5(t *testing.T) {
	os.Args = []string{"program5", "2015-12-09"}
	main()
}

func TestParse1(t *testing.T) {
	y, m, d, err := parse("2015-12-09")
	if err != nil {
		t.Errorf("Expected error: %v", err)
		return
	}
	if y != 2015 || m != 12 || d != 9 {
		t.Errorf("Got %d, %d, %d, want %d, %d, %d", 2015, 12, 9, y, m, d)
	}
}

func TestParse2(t *testing.T) {
	y, m, d, err := parse("2015-12-xx")
	if err == nil {
		t.Errorf("Should have returned non-nil error, got %d, %d, %d, %v instead.", y, m, d, err)
	}
}
