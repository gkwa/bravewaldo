package core11

import (
	"testing"
)

func TestReplaceURLsWithFriendlyNames(t *testing.T) {
	input := `
# asdf
this and that https://example.com test
[Google](https://example.com)
friday [Google](https://example.com)
`
	expected := `
# asdf
this and that [sample website](https://example.com) test
[Google](https://example.com)
friday [Google](https://example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestReplaceURLsWithFriendlyNamesCaseInsensitive(t *testing.T) {
	input := `
# asdf
this and that https://Example.com test
[Google](https://google.com)
friday [Google](https://example.com)
`
	expected := `
# asdf
this and that [sample website](https://Example.com) test
[Google](https://google.com)
friday [Google](https://example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
