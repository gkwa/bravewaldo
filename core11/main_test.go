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

type testCase struct {
	name           string
	inputText      string
	expectedOutput string
	options        ProcessOptions
}

func TestProcessMarkdown(t *testing.T) {
	testCases := []testCase{
		{
			name: "Replace URLs with friendly names",
			inputText: `
# asdf
this and that https://example.com test
[Google](https://example.com)
friday [Google](https://example.com)
`,
			expectedOutput: `
# asdf
this and that [sample website](https://example.com) test
[Google](https://example.com)
friday [Google](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Replace URLs with friendly names (case-insensitive)",
			inputText: `
# asdf
this and that https://Example.com test
[Google](https://google.com)
friday [Google](https://example.com)
`,
			expectedOutput: `
# asdf
this and that [sample website](https://Example.com) test
[Google](https://google.com)
friday [Google](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "URLs in parentheses should not be replaced",
			inputText: `
[ Google ]( https://google.com )
`,
			expectedOutput: `
[ Google ]( https://google.com )
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Multiple URLs on a single line should all be replaced",
			inputText: `
https://Example.com https://Example.com https://Example.com
`,
			expectedOutput: `
[sample website](https://Example.com) [sample website](https://Example.com) [sample website](https://Example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "URLs in markdown links should not be replaced",
			inputText: `
https://Example.com [test](https://Example.com) https://Example.com
`,
			expectedOutput: `
[sample website](https://Example.com) [test](https://Example.com) [sample website](https://Example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Existing markdown links should not be modified",
			inputText: `
[Example](https://Example.com)
`,
			expectedOutput: `
[Example](https://Example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "URLs with query parameters should be replaced",
			inputText: `
Check out https://example.com?param=value&another=123
`,
			expectedOutput: `
Check out [sample website2](https://example.com?param=value&another=123)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "URLs with fragments should be replaced",
			inputText: `
See https://example.com#section1 for more info
`,
			expectedOutput: `
See [sample website3](https://example.com#section1) for more info
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "URLs at the beginning and end of the line should be replaced",
			inputText: `
https://example.com is a great site and so is https://goOgle.com
`,
			expectedOutput: `
[sample website](https://example.com) is a great site and so is [search engine](https://goOgle.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Markdown links with double-quoted titles should be preserved without titles when IncludeTitle is false",
			inputText: `
Check out [Google](https://google.com "Search Engine")
And [Example](https://example.com "Sample Site")
`,
			expectedOutput: `
Check out [Google](https://google.com)
And [Example](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Markdown links with single-quoted titles should be preserved without titles when IncludeTitle is false",
			inputText: `
Check out [Google](https://google.com 'Search Engine')
And [Example](https://example.com 'Sample Site')
`,
			expectedOutput: `
Check out [Google](https://google.com)
And [Example](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Markdown links with titles and spaces should be preserved without titles when IncludeTitle is false",
			inputText: `
Visit [ Google ]( https://google.com  "Best Search Engine" )
Also [ Example ]( https://example.com  'Great Example Site' )
`,
			expectedOutput: `
Visit [Google](https://google.com)
Also [Example](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Mix of Markdown links with titles and plain URLs",
			inputText: `
[Google](https://google.com "Search") and https://example.com are great sites
Also check [Example](https://example.com 'Sample') and http://test.org
`,
			expectedOutput: `
[Google](https://google.com) and [sample website](https://example.com) are great sites
Also check [Example](https://example.com) and [testing site](http://test.org)
`,
			options: ProcessOptions{IncludeTitle: false},
		},
		{
			name: "Plain URLs with IncludeTitle true",
			inputText: `
Check out https://google.com and https://example.com
`,
			expectedOutput: `
Check out [search engine](https://google.com) and [sample website](https://example.com)
`,
			options: ProcessOptions{IncludeTitle: true},
		},
		{
			name: "Existing links with IncludeTitle true",
			inputText: `
[Google](https://google.com) and [Example](https://example.com "Sample")
Also [Test](https://test.org 'Testing')
`,
			expectedOutput: `
[Google](https://google.com) and [Example](https://example.com "Sample")
Also [Test](https://test.org "Testing")
`,
			options: ProcessOptions{IncludeTitle: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.inputText)
			var output bytes.Buffer
			err := ProcessMarkdown(input, &output, urlMap, tc.options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expectedOutput, output.String()); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}
		})
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

	var output bytes.Buffer
	options := ProcessOptions{IncludeTitle: true}
	err := ProcessMarkdown(strings.NewReader(input), &output, urlMap, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expected, output.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
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

	var output bytes.Buffer
	options := ProcessOptions{IncludeTitle: false}
	err := ProcessMarkdown(strings.NewReader(input), &output, urlMap, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expected, output.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}
