package core11

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/xurls/v2"
)

func ProcessMarkdown(input io.Reader, output io.Writer, urlMap map[string]string) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		urls := xurls.Strict().FindAllString(line, -1)

		for _, url := range urls {
			linkStart := strings.Index(line, "[")
			linkEnd := strings.Index(line, ")")

			if linkStart == -1 || linkEnd == -1 || linkStart > linkEnd {
				lowercaseURL := strings.ToLower(url)
				friendlyName, ok := urlMap[lowercaseURL]
				if ok {
					line = strings.Replace(line, url, fmt.Sprintf("[%s](%s)", friendlyName, url), 1)
				}
			}

		}

		fmt.Fprintln(output, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

func Main() error {
	urlMap := map[string]string{
		"https://example.com":          "sample website",
		"https://google.com":           "search engine",
		"http://test.org":              "testing site",
		"https://github.com/user/repo": "code repository",
	}

	inputFile, err := os.Open("testdata/input.md")
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("testdata/output.md")
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	return ProcessMarkdown(inputFile, outputFile, urlMap)
}
