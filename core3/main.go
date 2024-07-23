package core3

import (
	"bytes"
	"fmt"
	"log"
	"os"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
)

func Main() {
	filename := "testdata/input.md"

	source, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Create goldmark converter with markdown renderer object
	// Can pass functional Options as arguments. This example converts headings to ATX style.
	renderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	md := goldmark.New(goldmark.WithRenderer(renderer))

	// "Convert" markdown to formatted markdown
	buf := bytes.Buffer{}
	err = md.Convert(source, &buf)
	if err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	fmt.Println(buf.String())
}
