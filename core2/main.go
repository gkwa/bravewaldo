package core2

import (
	"bytes"
	"fmt"
	"log"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func Main() {
	// Create goldmark converter with markdown renderer object
	// Can pass functional Options as arguments. This example converts headings to ATX style.
	renderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	md := goldmark.New(
		goldmark.WithRenderer(renderer),
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// "Convert" markdown to formatted markdown
	source := `
My Document Title
=================

## Section 1

This is a paragraph.

* List item 1
* List item 2

### Subsection

More content here.
`
	buf := bytes.Buffer{}
	err := md.Convert([]byte(source), &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(buf.String())
}
