//go:build !skip_flaky
// +build !skip_flaky

package core11

import (
	"testing"
)

func TestMarkdownLinksWithTitlesAndSpaces(t *testing.T) {
	input := `
Visit [ Google ]( https://google.com  "Best Search Engine" )
Also [ Example ]( https://example.com  'Great Example Site' )
`
	expected := `
Visit [Google](https://google.com)
Also [Example](https://example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
