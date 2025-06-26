package main

import (
	"flag"
	"os"

	"github.com/rwirdemann/markdownsift"
)

func main() {
	path := flag.String("path", "/Users/ralfwirdemann/Documents/zettelkasten", "source directory to parse")
	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}
	markdownsift.WriteSnippets(os.Stdout, markdownsift.CollectSnippets(*path))
}
