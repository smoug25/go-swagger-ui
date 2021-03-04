package server

import (
	"net/http"
	"testing"
)

func TestGetFile(t *testing.T) {
	f, err := getFile("index.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(f) == 0 {
		t.Fatal("empty file")
	}
}

func TestExamole(t *testing.T) {
	// Set swagger file local path or url
	SetSwaggerFile("swagger/swagger.json")
	// add swagger ui handler
	http.HandleFunc("/", Serv)
	// start you http server
	http.ListenAndServe(":8080", nil)
}
