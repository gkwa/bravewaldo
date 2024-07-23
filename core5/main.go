package core5

import (
	"bytes"
	"log"
	"os"
	"regexp"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"mvdan.cc/xurls/v2"
)

var (
	inputFilename  = "testdata/input.md"
	outputFilename = "testdata/output.md"
)

func processURLs(input []byte) []byte {
	rxStrict := xurls.Strict()
	mdLinkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	protected := mdLinkRegex.ReplaceAllFunc(input, func(match []byte) []byte {
		return bytes.Replace(match, []byte("("), []byte("(protected:"), 1)
	})
	processed := rxStrict.ReplaceAllFunc(protected, func(match []byte) []byte {
		return match
	})
	return mdLinkRegex.ReplaceAllFunc(processed, func(match []byte) []byte {
		return bytes.Replace(match, []byte("(protected:"), []byte("("), 1)
	})
}

func Main() {
	source, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}
	processedSource := processURLs(source)
	md := goldmark.New(
		goldmark.WithRenderer(markdown.NewRenderer()),
	)
	var buf bytes.Buffer
	if err := md.Convert(processedSource, &buf); err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}
	output := buf.String()
	if err := os.WriteFile(outputFilename, []byte(output), 0o644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
}
