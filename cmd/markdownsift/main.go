package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rwirdemann/markdownsift"
)

func main() {
	path := flag.String("path", "/Users/ralfwirdemann/Documents/zettelkasten", "source directory to parse")
	flag.Parse()

	files, err := markdownsift.ListFiles(*path, "^\\d{4}-\\d{2}-\\d{2}\\.md$")
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		return
	}

	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)
	}

	file, err := os.Open("/Users/ralfwirdemann/Documents/zettelkasten/2025-06-20.md")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	result := markdownsift.CollectHashtaggedContent(file)

	for tag, blocks := range result {
		fmt.Printf("%s:\n", tag)
		for i, block := range blocks {
			fmt.Printf("Block %d:\n%s\n\n", i+1, block)
		}
	}
}
