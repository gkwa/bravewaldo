package core

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-logr/logr"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type RegexpLinkTransformer struct {
	LinkPattern *regexp.Regexp
	ReplUrl     []byte
}

func (t *RegexpLinkTransformer) LinkifyText(node *ast.Text, source []byte) {
}

func (t *RegexpLinkTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
}

func Example(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Example function")
	logger.Info("Processing markdown file")

	source, err := os.ReadFile("testdata/input.md")
	if err != nil {
		logger.Error(err, "Failed to read input file")
		return
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(source))

	var urls []string
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if autoLink, ok := n.(*ast.AutoLink); ok {
			url := string(autoLink.URL(source))
			urls = append(urls, url)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		logger.Error(err, "Error walking through the AST")
		return
	}

	logger.Info("Found URLs", "count", len(urls))
	for i, url := range urls {
		logger.Info(fmt.Sprintf("URL %d: %s", i+1, url))
	}

	logger.V(1).Info("Debug: Exiting Example function")
}
