package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwirdemann/markdownsift"
)

func main() {
	path := flag.String("path", "/Users/ralfwirdemann/Documents/zettelkasten", "source directory to parse")
	tags := flag.String("tags", "", "comma separated list of tags to be parsed (hash sign omitted, default = empty is all)")
	output := flag.String("output", "stdout", "output destination. values (stdout, file)")
	outputDir := flag.String("output-dir", "/Users/ralfwirdemann/Documents/zettelkasten/topics", "output directory for file output")

	flag.Parse()
	if err := validateFlags(*path, *output, *outputDir); err != nil {
		log.Println("Error:", err)
		flag.Usage()
		os.Exit(1)
	}

	// Prepend tags with hash sign
	var tt []string
	if *tags != "" {
		tt = strings.Split(*tags, ",")
		for i, tag := range tt {
			tt[i] = "#" + tag
		}
	}

	snippets := markdownsift.Filter(markdownsift.CollectSnippets(*path), tt)
	for tag, blocks := range snippets {
		writeBlocks(tag, blocks, *output, *outputDir)
	}
}

func validateFlags(path, output, outputDir string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if output != "stdout" && output != "file" {
		return fmt.Errorf("output must be either 'stdout' or 'file'")
	}
	if output == "file" && outputDir == "" {
		return fmt.Errorf("output-dir is required when output is 'file'")
	}
	return nil
}

func writeBlocks(tag string, blocks []markdownsift.Block, output, outputDir string) {
	writer := os.Stdout
	if output == "file" {
		// Create output directory if it doesn't exist
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		path := filepath.Join(outputDir, tag[1:]+".md")
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		writer = file
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	markdownsift.Write(writer, tag, blocks)
}
