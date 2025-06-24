package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	var snippets = map[string][]string{}
	for _, file := range files {
		func() {
			fmt.Printf("Processing file: %s\n", file)
			file, err := os.Open(filepath.Join(*path, file))
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			result := markdownsift.CollectHashtaggedContent(file)
			for tag, blocks := range result {
				snippets[tag] = append(snippets[tag], blocks...)
			}
		}()
	}

	for tag, blocks := range snippets {
		fmt.Printf("%s:\n", tag)
		for i, block := range blocks {
			fmt.Printf("Block %d:\n%s\n\n", i+1, block)
		}
	}
}
