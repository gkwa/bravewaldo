package core4

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/go-logr/logr"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func NewURLWrapperRenderer(logger logr.Logger) renderer.Renderer {
	logger.V(1).Info("Creating new URLWrapperRenderer")
	r := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	r.AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(URLWrapperNodeRenderer{logger: logger}, 100),
	))
	return r
}

type URLWrapperNodeRenderer struct {
	logger logr.Logger
}

func (r URLWrapperNodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	r.logger.V(1).Info("Registering renderAutoLink function")
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
}

func (r URLWrapperNodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.logger.V(1).Info("Entering renderAutoLink")
	if entering {
		n := node.(*ast.AutoLink)
		url := n.URL(source)
		wrappedURL := fmt.Sprintf("|%s|", url)
		r.logger.V(1).Info("Wrapping URL", "original", string(url), "wrapped", wrappedURL)
		_, err := w.WriteString(wrappedURL)
		if err != nil {
			r.logger.Error(err, "Failed to write wrapped URL")
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func Main(logger logr.Logger) {
	logger.V(1).Info("Entering Main function")
	filename := "testdata/input.md"

	source, err := os.ReadFile(filename)
	if err != nil {
		logger.Error(err, "Error reading file")
		log.Fatalf("Error reading file: %v", err)
	}

	logger.V(1).Info("Creating new Goldmark instance")
	md := goldmark.New(
		goldmark.WithRenderer(NewURLWrapperRenderer(logger)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var buf bytes.Buffer
	logger.V(1).Info("Converting markdown")
	if err := md.Convert(source, &buf); err != nil {
		logger.Error(err, "Error converting markdown")
		log.Fatalf("Error converting markdown: %v", err)
	}

	logger.V(1).Info("Conversion complete", "output", buf.String())
	fmt.Println(buf.String())
}
