package core11

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
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
		},
		{
			name: "URLs in parentheses should not be replaced",
			inputText: `
[ Google ]( https://google.com )
`,
			expectedOutput: `
[ Google ]( https://google.com )
`,
		},
		{
			name: "Multiple URLs on a single line should all be replaced",
			inputText: `
https://Example.com https://Example.com https://Example.com
`,
			expectedOutput: `
[sample website](https://Example.com) [sample website](https://Example.com) [sample website](https://Example.com)
`,
		},
		{
			name: "URLs in markdown links should not be replaced",
			inputText: `
https://Example.com [test](https://Example.com) https://Example.com
`,
			expectedOutput: `
[sample website](https://Example.com) [test](https://Example.com) [sample website](https://Example.com)
`,
		},
		{
			name: "Existing markdown links should not be modified",
			inputText: `
[Example](https://Example.com)
`,
			expectedOutput: `
[Example](https://Example.com)
`,
		},
		{
			name: "URLs with query parameters should be replaced",
			inputText: `
Check out https://example.com?param=value&another=123
`,
			expectedOutput: `
Check out [sample website2](https://example.com?param=value&another=123)
`,
		},
		{
			name: "URLs with fragments should be replaced",
			inputText: `
See https://example.com#section1 for more info
`,
			expectedOutput: `
See [sample website3](https://example.com#section1) for more info
`,
		},
		{
			name: "URLs at the beginning and end of the line should be replaced",
			inputText: `
https://example.com is a great site and so is https://goOgle.com
`,
			expectedOutput: `
[sample website](https://example.com) is a great site and so is [search engine](https://goOgle.com)
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.inputText)
			var output bytes.Buffer
			err := ProcessMarkdown(input, &output, urlMap)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expectedOutput, output.String()); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFuzzProcessMarkdown(t *testing.T) {
	f := fuzz.New()
	var input string

	for i := 0; i < 1000; i++ {
		f.Fuzz(&input)

		inputReader := strings.NewReader(input)
		var output bytes.Buffer
		err := ProcessMarkdown(inputReader, &output, urlMap)
		if err != nil {
			t.Errorf("ProcessMarkdown failed on fuzz input: %v\nInput: %s", err, input)
		}
	}
}

func TestRandomizedInput(t *testing.T) {
	urlKeys := make([]string, 0, len(urlMap))
	for k := range urlMap {
		urlKeys = append(urlKeys, k)
	}

	for i := 0; i < 100; i++ {
		input := generateRandomInput(urlKeys)
		inputReader := strings.NewReader(input)
		var output bytes.Buffer
		err := ProcessMarkdown(inputReader, &output, urlMap)
		if err != nil {
			t.Errorf("ProcessMarkdown failed on random input: %v\nInput: %s", err, input)
		}

		t.Logf("Random Test %d:\nInput: %s\nOutput: %s\n", i+1, input, output.String())
	}
}

func generateRandomInput(urls []string) string {
	words := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
	var result strings.Builder

	for i := 0; i < rand.Intn(20)+1; i++ {
		if rand.Float32() < 0.3 && len(urls) > 0 {
			url := urls[rand.Intn(len(urls))]
			if rand.Float32() < 0.5 {
				result.WriteString(fmt.Sprintf("[%s](%s) ", words[rand.Intn(len(words))], url))
			} else {
				result.WriteString(url + " ")
			}
		} else {
			result.WriteString(words[rand.Intn(len(words))] + " ")
		}
	}

	return strings.TrimSpace(result.String())
}
