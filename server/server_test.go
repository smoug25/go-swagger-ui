package server

import "testing"

func TestGetFile(t *testing.T) {
	f, err := getFile("index.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(f) == 0 {
		t.Fatal("empty file")
	}
}
