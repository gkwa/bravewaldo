package core

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func NewURLWrapperRenderer(logger logr.Logger) renderer.Renderer {
	logger.V(1).Info("Creating new URLWrapperRenderer")
	r := renderer.NewRenderer()
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

func Example(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Example function")

	source, err := os.ReadFile("testdata/input.md")
	if err != nil {
		logger.Error(err, "Failed to read input file")
		return
	}

	logger.V(1).Info("Creating new Goldmark instance")
	md := goldmark.New(
		goldmark.WithRenderer(NewURLWrapperRenderer(logger)),
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList, meta.Meta),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	initialAST, err := dumpAST(doc)
	if err != nil {
		logger.Error(err, "Failed to dump initial AST")
		return
	}
	logger.V(1).Info("Initial AST structure", "structure", initialAST)

	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		logger.V(1).Info("Walking node", "type", fmt.Sprintf("%T", n), "kind", n.Kind())

		if autoLink, ok := n.(*ast.AutoLink); ok {
			url := autoLink.URL(source)
			wrappedUrl := fmt.Sprintf("|%s|", url)
			logger.V(1).Info("Found AutoLink", "url", string(url), "wrapped", wrappedUrl)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		logger.Error(err, "Error walking through the AST")
		return
	}

	finalAST, err := dumpAST(doc)
	if err != nil {
		logger.Error(err, "Failed to dump final AST")
		return
	}
	logger.V(1).Info("Final AST structure", "structure", finalAST)

	var buf bytes.Buffer
	logger.V(1).Info("Starting markdown rendering")
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		logger.Error(err, "Error rendering markdown")
		return
	}
	logger.V(1).Info("Finished markdown rendering")

	output := buf.String()
	logger.V(1).Info("Rendered output", "output", output)

	if err := os.WriteFile("testdata/output.md", buf.Bytes(), 0o644); err != nil {
		logger.Error(err, "Error writing output file")
		return
	}

	logger.V(1).Info("Debug: Exiting Example function")
}

func dumpAST(n ast.Node) (string, error) {
	var buf strings.Builder
	level := 0
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			fmt.Fprintf(&buf, "%s%s {\n", strings.Repeat("  ", level), n.Kind())
			level++
		} else {
			level--
			fmt.Fprintf(&buf, "%s}\n", strings.Repeat("  ", level))
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", fmt.Errorf("error walking AST: %w", err)
	}
	return buf.String(), nil
}
