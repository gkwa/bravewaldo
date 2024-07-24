package core10

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
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var (
	inputFilename  = "testdata/input.md"
	outputFilename = "testdata/output.md"
)

var urlMap = map[string]string{
	"https://example.com":          "sample website",
	"https://google.com":           "search engine",
	"http://test.org":              "testing site",
	"https://github.com/user/repo": "code repository",
}

func newURLRewriteRenderer(logger logr.Logger) renderer.Renderer {
	logger.V(1).Info("Creating new URLRewriteRenderer")
	r := markdown.NewRenderer()
	r.AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(urlRewriteNodeRenderer{logger: logger}, 1),
	))
	return r
}

type urlRewriteNodeRenderer struct {
	logger logr.Logger
}

func (r urlRewriteNodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	r.logger.V(1).Info("Registering renderAutoLink function")
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindText, r.renderText)
}

func (r urlRewriteNodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.logger.V(1).Info("Entering renderAutoLink")
	if entering {
		n := node.(*ast.AutoLink)
		url := string(n.URL(source))
		r.logger.V(1).Info("Processing URL", "url", url)
		if value, ok := urlMap[url]; ok {
			r.logger.V(1).Info("Rewriting AutoLink", "url", url, "value", value)
			fmt.Fprintf(w, "[%s](%s)", value, url)
			return ast.WalkSkipChildren, nil
		}
		r.logger.V(1).Info("URL not found in map, leaving as is", "url", url)
	}
	return ast.WalkContinue, nil
}

func (r urlRewriteNodeRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.logger.V(1).Info("Entering renderText")
	if entering {
		n := node.(*ast.Text)
		segment := n.Segment
		value := segment.Value(source)
		r.logger.V(1).Info("Processing Text", "value", string(value))
	}
	return ast.WalkContinue, nil
}

func Main(logger logr.Logger) {
	logger.V(1).Info("Entering Main function")
	source, err := os.ReadFile(inputFilename)
	if err != nil {
		logger.Error(err, "Error reading input file")
		log.Fatalf("Error reading input file: %v", err)
	}

	logger.V(1).Info("Creating new Goldmark instance")
	md := goldmark.New(
		goldmark.WithRenderer(newURLRewriteRenderer(logger)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	logger.V(1).Info("Parsing markdown")
	doc := md.Parser().Parse(text.NewReader(source))

	logger.V(1).Info("Walking AST")
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			logger.V(1).Info("Node", "type", fmt.Sprintf("%T", n), "kind", n.Kind())
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		logger.Error(err, "Error walking AST")
	}

	var buf bytes.Buffer
	logger.V(1).Info("Rendering markdown")
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		logger.Error(err, "Error rendering markdown")
		log.Fatalf("Error rendering markdown: %v", err)
	}

	logger.V(1).Info("Writing output file")
	if err := os.WriteFile(outputFilename, buf.Bytes(), 0o644); err != nil {
		logger.Error(err, "Error writing output file")
		log.Fatalf("Error writing output file: %v", err)
	}

	logger.V(1).Info("Main function completed successfully")
}
