package core11

import (
	"testing"
)

func TestMarkdownLinksWithDoubleQuotedTitles(t *testing.T) {
	input := `
Check out [Google](https://google.com "Search Engine")
And [Example](https://example.com "Sample Site")
`
	expected := `
Check out [Google](https://google.com)
And [Example](https://example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestMarkdownLinksWithSingleQuotedTitles(t *testing.T) {
	input := `
Check out [Google](https://google.com 'Search Engine')
And [Example](https://example.com 'Sample Site')
`
	expected := `
Check out [Google](https://google.com)
And [Example](https://example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
