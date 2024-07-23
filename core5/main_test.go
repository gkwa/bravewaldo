package core5

import (
	"bytes"
	"os"
	"testing"
)

func TestProcessURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No URLs",
			input:    "This is a test with no URLs",
			expected: "This is a test with no URLs",
		},
		{
			name:     "Single URL",
			input:    "Check out https://example.com for more info",
			expected: "Check out https://example.com for more info",
		},
		{
			name:     "Multiple URLs",
			input:    "Visit https://example.com and http://test.org",
			expected: "Visit https://example.com and http://test.org",
		},
		{
			name:     "Markdown link",
			input:    "Click [here](https://example.com) for more",
			expected: "Click [here](https://example.com) for more",
		},
		{
			name:     "Mixed content",
			input:    "Check [this](https://example.com) and https://test.org",
			expected: "Check [this](https://example.com) and https://test.org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processURLs([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("processURLs() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestMain(t *testing.T) {
	inputFile, err := os.CreateTemp("", "input*.md")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	outputFile, err := os.CreateTemp("", "output*.md")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outputFile.Name())

	testInput := `# Test Markdown
This is a [test link](https://example.com).
Here's a plain URL: https://test.org
`
	if _, err := inputFile.Write([]byte(testInput)); err != nil {
		t.Fatalf("Failed to write to temp input file: %v", err)
	}
	inputFile.Close()

	oldInputFilename, oldOutputFilename := inputFilename, outputFilename
	inputFilename, outputFilename = inputFile.Name(), outputFile.Name()
	defer func() {
		inputFilename, outputFilename = oldInputFilename, oldOutputFilename
	}()

	Main()

	output, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !bytes.Equal(output, []byte(testInput)) {
		t.Errorf("Main() output = %v, want %v", string(output), testInput)
	}
}

func TestSourceOutputComparison(t *testing.T) {
	inputFile, err := os.CreateTemp("", "input*.md")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	outputFile, err := os.CreateTemp("", "output*.md")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outputFile.Name())

	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "Mixed content",
			input: `# Sample Markdown
This is a [Markdown link](https://example.com).
Here's a plain URL: https://test.org
And here's some code:
` + "```" + `
https://code-example.com
` + "```" + `
`,
		},
		{
			name: "Only plain text",
			input: `Just some plain text.
No URLs or Markdown links here.
`,
		},
		{
			name: "Multiple Markdown links",
			input: `Check out [this link](https://example1.com) and [that link](https://example2.com).
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := inputFile.Write([]byte(tc.input)); err != nil {
				t.Fatalf("Failed to write to temp input file: %v", err)
			}
			inputFile.Close()

			oldInputFilename, oldOutputFilename := inputFilename, outputFilename
			inputFilename, outputFilename = inputFile.Name(), outputFile.Name()
			defer func() {
				inputFilename, outputFilename = oldInputFilename, oldOutputFilename
			}()

			Main()

			output, err := os.ReadFile(outputFile.Name())
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if !bytes.Equal(output, []byte(tc.input)) {
				t.Errorf("Output doesn't match input.\nInput:\n%s\nOutput:\n%s", tc.input, string(output))
			}

			if err := inputFile.Truncate(0); err != nil {
				t.Fatalf("Failed to truncate input file: %v", err)
			}
			if _, err := inputFile.Seek(0, 0); err != nil {
				t.Fatalf("Failed to seek input file: %v", err)
			}
		})
	}
}
