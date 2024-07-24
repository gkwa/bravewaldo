package core11

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var urlMap = map[string]string{
	"https://example.com":                         "sample website",
	"https://google.com":                          "search engine",
	"http://test.org":                             "testing site",
	"https://github.com/user/repo":                "code repository",
	"https://example.com?param=value&another=123": "sample website2",
	"https://example.com#section1":                "sample website3",
}

func testProcessMarkdown(t *testing.T, input, expected string, options ProcessOptions) {
	t.Helper()
	var output bytes.Buffer
	err := ProcessMarkdown(strings.NewReader(input), &output, urlMap, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expected, output.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

func TestProcessMarkdownWithTitles(t *testing.T) {
	input := `
Check out [Google](https://google.com "Search Engine")
And [Example](https://example.com 'Sample Site')
Also https://example.com
`
	expected := `
Check out [Google](https://google.com "Search Engine")
And [Example](https://example.com "Sample Site")
Also [sample website](https://example.com)
`

	options := ProcessOptions{IncludeTitle: true}
	testProcessMarkdown(t, input, expected, options)
}

func TestProcessMarkdownWithoutTitles(t *testing.T) {
	input := `
Check out [Google](https://google.com "Search Engine")
And [Example](https://example.com 'Sample Site')
Also https://example.com
`
	expected := `
Check out [Google](https://google.com)
And [Example](https://example.com)
Also [sample website](https://example.com)
`

	options := ProcessOptions{IncludeTitle: false}
	testProcessMarkdown(t, input, expected, options)
}
