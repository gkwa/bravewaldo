package core9

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

	renderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	md := goldmark.New(goldmark.WithRenderer(renderer))

	buf := bytes.Buffer{}
	err = md.Convert(source, &buf)
	if err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	fmt.Println(buf.String())
}
