package core

import (
	"bytes"
	"fmt"
	"io"
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
)

type URLWrapperRenderer struct {
	renderer.Renderer
}

func NewURLWrapperRenderer() renderer.Renderer {
	return &URLWrapperRenderer{
		Renderer: renderer.NewRenderer(),
	}
}

func (r *URLWrapperRenderer) Render(w io.Writer, source []byte, n ast.Node) error {
	return ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if autoLink, ok := n.(*ast.AutoLink); ok {
				url := autoLink.URL(source)
				wrappedURL := fmt.Sprintf("|%s|", url)
				_, err := w.Write([]byte(wrappedURL))
				if err != nil {
					return ast.WalkStop, err
				}
				return ast.WalkSkipChildren, nil
			}
		}
		err := r.Renderer.Render(w, source, n)
		if err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkContinue, nil
	})
}

func (r *URLWrapperRenderer) AddOptions(...renderer.Option) {}

func Example(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Example function")
	logger.Info("Processing markdown file")

	source, err := os.ReadFile("testdata/input.md")
	if err != nil {
		logger.Error(err, "Failed to read input file")
		return
	}

	md := goldmark.New(
		goldmark.WithRenderer(NewURLWrapperRenderer()),
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

	var urls []string
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		logger.V(1).Info("Walking node", "type", fmt.Sprintf("%T", n), "kind", n.Kind())

		if autoLink, ok := n.(*ast.AutoLink); ok {
			url := autoLink.URL(source)
			wrappedUrl := fmt.Sprintf("|%s|", url)
			urls = append(urls, wrappedUrl)
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

	logger.Info("Found URLs", "count", len(urls))
	for i, url := range urls {
		logger.Info(fmt.Sprintf("URL %d: %s", i+1, url))
	}

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