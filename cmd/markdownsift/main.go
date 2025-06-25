package main

import (
	"flag"
	"fmt"

	"github.com/rwirdemann/markdownsift"
)

func main() {
	path := flag.String("path", "/Users/ralfwirdemann/Documents/zettelkasten", "source directory to parse")
	flag.Parse()

	snippets := markdownsift.CollectSnippets(*path)
	for tag, blocks := range snippets {
		fmt.Printf("%s:\n", tag)
		for i, block := range blocks {
			fmt.Printf("Block %d:\n%s\n\n", i+1, block)
		}
	}
}
