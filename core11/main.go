package core11

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func ProcessMarkdown(input io.Reader, output io.Writer, urlMap map[string]string) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)|\bhttps?://\S+`)
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			if strings.HasPrefix(match, "[") {
				return match
			}
			url := match
			lowercaseURL := strings.ToLower(url)
			friendlyName, ok := urlMap[lowercaseURL]
			if ok {
				return fmt.Sprintf("[%s](%s)", friendlyName, url)
			}
			return url
		})
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
