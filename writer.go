package markdownsift

import (
	"fmt"
	"log"
)

// Writer provides methods to write snippet data to an output target like stdout or a file.
type Writer interface {
	Create(name string) error
	Write(p []byte) (n int, err error)
	Close() error
}

func WriteSnippets(snippets map[string][]Block, writer Writer) {
	for tag, blocks := range snippets {
		func() {
			if err := writer.Create(tag); err != nil {
				log.Fatalf(err.Error())
			}
			defer func(writer Writer) {
				_ = writer.Close()
			}(writer)

			if _, err := fmt.Fprintf(writer, "# Content tagged by %s\n\n", tag); err != nil {
				log.Fatalf(err.Error())
			}
			for _, block := range blocks {
				if _, err := fmt.Fprintf(writer, "[%s](../%s.md):\n%s\n\n", block.Date.Format(dateFormat), block.Date.Format(dateFormat), block.Content); err != nil {
					log.Fatalf(err.Error())
				}
			}
		}()
	}
}
