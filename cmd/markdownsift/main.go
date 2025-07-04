package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwirdemann/markdownsift"
	"github.com/rwirdemann/markdownsift/file"
	_os "github.com/rwirdemann/markdownsift/os"
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
	var writer markdownsift.Writer
	var err error
	switch *output {
	case "stdout":
		writer = _os.NewWriter()
	case "file":
		writer, err = file.NewWriter(*outputDir)
		if err != nil {
			log.Println("Error:", err)
			flag.Usage()
			os.Exit(1)
		}
	}
	markdownsift.WriteSnippets(snippets, writer)
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
