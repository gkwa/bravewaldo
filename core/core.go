package core

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-logr/logr"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type RegexpLinkTransformer struct {
	LinkPattern *regexp.Regexp
	ReplUrl     []byte
}

func (t *RegexpLinkTransformer) LinkifyText(node *ast.Text, source []byte) {
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

	link := ast.NewLink()
	link.AppendChild(link, ast.NewTextSegment(lSegment))
	link.Destination = t.LinkPattern.ReplaceAll(lSegment.Value(source), t.ReplUrl)
	parent.InsertBefore(parent, node, link)

	node.Segment = tSegment.WithStart(lSegment.Stop)

	if node.Segment.Len() > 0 {
		t.LinkifyText(node, source)
	}
}

func (t *RegexpLinkTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if node.Kind() == ast.KindLink || node.Kind() == ast.KindAutoLink {
			return ast.WalkSkipChildren, nil
		}
		if node.Kind() == ast.KindText {
			textNode := node.(*ast.Text)
			t.LinkifyText(textNode, source)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		log.Fatal("Error encountered while transforming AST:", err)
	}
}

func Example(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Example function")
	logger.Info("Processing markdown file")

	source, err := os.ReadFile("testdata/input.md")
	if err != nil {
		logger.Error(err, "Failed to read input file")
		return
	}

	transformer := RegexpLinkTransformer{
		LinkPattern: regexp.MustCompile(`TICKET-\d+`),
		ReplUrl:     []byte("https://example.com/TICKET?query=$0"),
	}
	prioritizedTransformer := util.Prioritized(&transformer, 0)
	gm := goldmark.New(
		goldmark.WithRenderer(markdown.NewRenderer()),
		goldmark.WithParserOptions(parser.WithASTTransformers(prioritizedTransformer)),
	)
	buf := bytes.Buffer{}

	err = gm.Convert(source, &buf)
	if err != nil {
		logger.Error(err, "Encountered Markdown conversion error")
		return
	}
	fmt.Print(buf.String())

	logger.V(1).Info("Debug: Exiting Example function")
}
