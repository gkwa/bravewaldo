package core5

import (
	"bytes"
	"fmt"
	"log"
	"os"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"mvdan.cc/xurls/v2"
)

func wrapURLs(input []byte) []byte {
	rxStrict := xurls.Strict()
	return rxStrict.ReplaceAllFunc(input, func(match []byte) []byte {
		return []byte(fmt.Sprintf("|%s|", match))
	})
}

func Main() {
	inputFilename := "testdata/input.md"
	outputFilename := "testdata/output.md"

	source, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	wrappedSource := wrapURLs(source)

	md := goldmark.New(
		goldmark.WithRenderer(markdown.NewRenderer()),
	)

	var buf bytes.Buffer
	if err := md.Convert(wrappedSource, &buf); err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	output := buf.String()
	if err := os.WriteFile(outputFilename, []byte(output), 0o644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
}
