package core4

import (
	"bytes"
	"fmt"
	"log"
	"os"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type URLWrapperRenderer struct {
	markdown.Renderer
}

func NewURLWrapperRenderer() renderer.Renderer {
	r := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	r.AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(URLWrapperNodeRenderer{}, 100),
	))
	return r
}

type URLWrapperNodeRenderer struct{}

func (r URLWrapperNodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
}

func (r URLWrapperNodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.AutoLink)
		url := n.URL(source)
		wrappedURL := fmt.Sprintf("|%s|", url)
		_, err := w.WriteString(wrappedURL)
		if err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func Main() {
	filename := "testdata/input.md"

	source, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	md := goldmark.New(
		goldmark.WithRenderer(NewURLWrapperRenderer()),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	fmt.Println(buf.String())
}
