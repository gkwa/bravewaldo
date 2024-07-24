package core7

import (
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func Main() {
	inputFile := "testdata/input.md"

	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	var autoLinkURLs []string
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if autoLink, ok := n.(*ast.AutoLink); ok {
				url := string(autoLink.URL(source))
				autoLinkURLs = append(autoLinkURLs, url)
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		fmt.Printf("Error walking AST: %v\n", err)
		return
	}

	fmt.Println("Found AutoLink URLs:")
	for _, url := range autoLinkURLs {
		fmt.Println(url)
	}
}
