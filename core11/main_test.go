package core11

import (
	"bytes"
	"strings"
	"testing"
)

type testCase struct {
	name           string
	inputText      string
	expectedOutput string
	urlMap         map[string]string
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
			urlMap: map[string]string{
				"https://example.com":          "sample website",
				"https://google.com":           "search engine",
				"http://test.org":              "testing site",
				"https://github.com/user/repo": "code repository",
			},
		},

		{
			name: "Replace URLs with friendly names",
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
			urlMap: map[string]string{
				"https://example.com":          "sample website",
				"https://google.com":           "search engine",
				"http://test.org":              "testing site",
				"https://github.com/user/repo": "code repository",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.inputText)
			var output bytes.Buffer
			err := ProcessMarkdown(input, &output, tc.urlMap)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output.String() != tc.expectedOutput {
				t.Errorf("output does not match expected output\nGot:\n%s\nExpected:\n%s", output.String(), tc.expectedOutput)
			}
		})
	}
}
