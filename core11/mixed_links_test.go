package core11

import (
	"testing"
)

func TestMixOfMarkdownLinksAndPlainURLs(t *testing.T) {
	input := `
[Google](https://google.com "Search") and https://example.com are great sites
Also check [Example](https://example.com 'Sample') and http://test.org
`
	expected := `
[Google](https://google.com) and [sample website](https://example.com) are great sites
Also check [Example](https://example.com) and [testing site](http://test.org)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestPlainURLsWithIncludeTitleTrue(t *testing.T) {
	input := `
Check out https://google.com and https://example.com
`
	expected := `
Check out [search engine](https://google.com) and [sample website](https://example.com)
`
	options := ProcessOptions{IncludeTitle: true}
	testProcessMarkdown(t, input, expected, options)
}

func TestExistingLinksWithIncludeTitleTrue(t *testing.T) {
	input := `
[Google](https://google.com) and [Example](https://example.com "Sample")
Also [Test](https://test.org 'Testing')
`
	expected := `
[Google](https://google.com) and [Example](https://example.com "Sample")
Also [Test](https://test.org "Testing")
`
	options := ProcessOptions{IncludeTitle: true}
	testProcessMarkdown(t, input, expected, options)
}
