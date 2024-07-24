package core8

import (
	"bytes"
	"fmt"
	"log"
	"os"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func Main() {
	renderer := markdown.NewRenderer()
	md := goldmark.New(
		goldmark.WithRenderer(renderer),
		goldmark.WithExtensions(extension.GFM, meta.Meta),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	source, err := os.ReadFile("testdata/input.md")
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	buf := bytes.Buffer{}
	err = md.Convert(source, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(buf.String())
}
