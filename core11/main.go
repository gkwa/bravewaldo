package core11

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type MarkdownLink struct {
	Name  string
	URL   string
	Title string
}

type ProcessOptions struct {
	IncludeTitle bool
}

func ProcessMarkdown(input io.Reader, output io.Writer, urlMap map[string]string, options ProcessOptions) error {
	scanner := bufio.NewScanner(input)

	re := regexp.MustCompile(`\[([^\]]+)\]\s*\(([^)]+)\)|\bhttps?://\S+`)

	for scanner.Scan() {
		line := scanner.Text()
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			link := parseMarkdownLink(match)
			if link.Name != "" && link.URL != "" {
				return formatMarkdownLink(link, options.IncludeTitle)
			}

			surroundingText := line
			index := strings.Index(surroundingText, match)
			if index > 0 && strings.HasSuffix(strings.TrimSpace(surroundingText[:index]), "[") &&
				index+len(match) < len(surroundingText) && strings.HasPrefix(strings.TrimSpace(surroundingText[index+len(match):]), ")") {
				return match
			}

			lowercaseURL := strings.ToLower(link.URL)
			friendlyName, ok := urlMap[lowercaseURL]
			if ok {
				link.Name = friendlyName
				return formatMarkdownLink(link, options.IncludeTitle)
			}
			return link.URL
		})

		fmt.Fprintln(output, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

func parseMarkdownLink(text string) MarkdownLink {
	re := regexp.MustCompile(`^\s*\[([^\]]+)\]\s*\(([^)\s]+)(?:\s+(?:"([^"]+)"|'([^']+)')?)?\s*\)\s*$`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		link := MarkdownLink{
			Name: strings.TrimSpace(matches[1]),
			URL:  strings.TrimSpace(matches[2]),
		}
		if len(matches) >= 4 {
			if matches[3] != "" {
				link.Title = matches[3] // Double-quoted title
			} else if matches[4] != "" {
				link.Title = matches[4] // Single-quoted title
			}
		}
		return link
	}
	return MarkdownLink{URL: text}
}

func formatMarkdownLink(link MarkdownLink, includeTitle bool) string {
	if includeTitle && link.Title != "" {
		if strings.Contains(link.Title, `"`) {
			return fmt.Sprintf("[%s](%s '%s')", link.Name, link.URL, link.Title)
		}
		return fmt.Sprintf("[%s](%s \"%s\")", link.Name, link.URL, link.Title)
	}
	return fmt.Sprintf("[%s](%s)", link.Name, link.URL)
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

	options := ProcessOptions{IncludeTitle: false}
	return ProcessMarkdown(inputFile, outputFile, urlMap, options)
}
