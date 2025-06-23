package main

import (
	"fmt"
	"os"

	"github.com/rwirdemann/markdownsift"
)

func main() {
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
