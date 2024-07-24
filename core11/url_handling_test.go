package core11

import (
	"testing"
)

func TestURLsInParentheses(t *testing.T) {
	input := `
[ Google ]( https://google.com )
`
	expected := `
[ Google ]( https://google.com )
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestMultipleURLsOnSingleLine(t *testing.T) {
	input := `
https://Example.com https://Example.com https://Example.com
`
	expected := `
[sample website](https://Example.com) [sample website](https://Example.com) [sample website](https://Example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestURLsInMarkdownLinks(t *testing.T) {
	input := `
https://Example.com [test](https://Example.com) https://Example.com
`
	expected := `
[sample website](https://Example.com) [test](https://Example.com) [sample website](https://Example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}

func TestExistingMarkdownLinks(t *testing.T) {
	input := `
[Example](https://Example.com)
`
	expected := `
[Example](https://Example.com)
`
	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
