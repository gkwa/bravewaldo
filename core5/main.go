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
	"mvdan.cc/xurls/v2"
)

type URLWrapperTransformer struct {
	urlRegex *regexp.Regexp
}

func (t *URLWrapperTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	if err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindText {
			return ast.WalkContinue, nil
		}

		textNode := n.(*ast.Text)
		segment := textNode.Segment
		original := segment.Value(source)

		wrapped := t.urlRegex.ReplaceAllFunc(original, func(m []byte) []byte {
			return []byte(fmt.Sprintf("|%s|", m))
		})

		if !bytes.Equal(original, wrapped) {
			newNode := ast.NewString(wrapped)
			newNode.SetRaw(textNode.IsRaw())
			n.Parent().ReplaceChild(n.Parent(), n, newNode)
		}

		return ast.WalkContinue, nil
	}); err != nil {
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
		urlRegex: xurls.Strict(),
	}

	renderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	md := goldmark.New(
		goldmark.WithRenderer(renderer),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(transformer, 100),
			),
		),
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
