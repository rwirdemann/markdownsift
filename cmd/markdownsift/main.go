package main

import (
	"flag"
	"os"
	"strings"

	"github.com/rwirdemann/markdownsift"
)

func main() {
	path := flag.String("path", "/Users/ralfwirdemann/Documents/zettelkasten", "source directory to parse")
	tags := flag.String("tags", "", "comma separated list of tags to be parsed (hash sign omitted, default = empty is all)")
	flag.Parse()
	if *path == "" {
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

	snippets := markdownsift.CollectSnippets(*path)
	markdownsift.WriteSnippets(os.Stdout, snippets, tt)
}
