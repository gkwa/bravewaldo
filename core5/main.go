package core5

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type URLWrapperTransformer struct {
	LinkPattern *regexp.Regexp
}

func (t *URLWrapperTransformer) WrapURLs(node *ast.Text, source []byte) {
	parent := node.Parent()
	tSegment := node.Segment
	match := t.LinkPattern.FindIndex(tSegment.Value(source))
	if match == nil {
		return
	}

	lSegment := text.NewSegment(tSegment.Start+match[0], tSegment.Start+match[1])

	if lSegment.Start != tSegment.Start {
		bText := ast.NewTextSegment(tSegment.WithStop(lSegment.Start))
		parent.InsertBefore(parent, node, bText)
	}

	wrappedURL := fmt.Sprintf("|%s|", lSegment.Value(source))
	wrappedNode := ast.NewString([]byte(wrappedURL))
	parent.InsertBefore(parent, node, wrappedNode)

	node.Segment = tSegment.WithStart(lSegment.Stop)

	if node.Segment.Len() > 0 {
		t.WrapURLs(node, source)
	}
}

func (t *URLWrapperTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if node.Kind() == ast.KindText {
			textNode := node.(*ast.Text)
			t.WrapURLs(textNode, source)
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		log.Printf("Error walking AST: %v", err)
	}
}

func Main() {
	inputFilename := "testdata/input.md"
	outputFilename := "testdata/output.md"

	source, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	transformer := &URLWrapperTransformer{
		LinkPattern: regexp.MustCompile(`https?://\S+`),
	}
	prioritizedTransformer := util.Prioritized(transformer, 0)

	md := goldmark.New(
		goldmark.WithRenderer(markdown.NewRenderer()),
		goldmark.WithParserOptions(parser.WithASTTransformers(prioritizedTransformer)),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	output := buf.String()
	if err := os.WriteFile(outputFilename, []byte(output), 0o644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Output written to %s\n", outputFilename)
}
