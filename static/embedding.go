// +build ignore

package main

import (
	"log"
	"net/http"
	"github.com/shurcooL/vfsgen"
)

func main() {
	log.Print("Start generate static embed.")
	err := vfsgen.Generate(
		http.Dir("../dist"),
		vfsgen.Options{
			Filename:     "./static_vfsdata.go",
			PackageName:  "static",
			VariableName: "EmbedStatic",
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Print("Generation finish success.")
}
