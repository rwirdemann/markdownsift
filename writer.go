package markdownsift

import (
	"fmt"
	"log"
)

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

			if _, err := fmt.Fprintf(writer, "# Content tagged by %s\n", tag); err != nil {
				log.Fatalf(err.Error())
			}
			for _, block := range blocks {
				if _, err := fmt.Fprintf(writer, "%s:\n%s\n\n", block.Date.Format(dateFormat), block.Content); err != nil {
					log.Fatalf(err.Error())
				}
			}
		}()
	}
}
