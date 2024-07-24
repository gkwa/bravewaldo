package core11

import (
	"testing"
)

func TestURLsWithQueryParameters(t *testing.T) {
	input := `
Check out https://example.com?param=value&another=123
`
	expected := `
Check out [sample website2](https://example.com?param=value&another=123)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestURLsWithFragments(t *testing.T) {
	input := `
See https://example.com#section1 for more info
`
	expected := `
See [sample website3](https://example.com#section1) for more info
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestURLsAtBeginningAndEnd(t *testing.T) {
	input := `
https://example.com is a great site and so is https://goOgle.com
`
	expected := `
[sample website](https://example.com) is a great site and so is [search engine](https://goOgle.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
