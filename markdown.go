package markdownsift

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// ListFiles returns a list of files in the given path that match the given pattern.
func ListFiles(path string, pattern string) ([]string, error) {
	var matchingFiles []string

	// Compile the regex pattern
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern '%s': %w", pattern, err)
	}

	// Read the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", path, err)
	}

	// Filter files that match the pattern
	for _, entry := range entries {
		if !entry.IsDir() && regex.MatchString(entry.Name()) {
			matchingFiles = append(matchingFiles, entry.Name())
		}
	}

	return matchingFiles, nil
}

// CollectHashtaggedContent returns a map of hashtags to their content blocks.
func CollectHashtaggedContent(reader io.Reader) map[string][]string {
	result := make(map[string][]string)

	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return result
	}

	lines := strings.Split(string(content), "\n")
	hashtagRegex := regexp.MustCompile(`#\w+`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		// Find all hashtags in the line
		hashtags := hashtagRegex.FindAllString(line, -1)

		if len(hashtags) > 0 {
			// This line contains hashtags, collect the block
			block := []string{line}

			// Collect following lines until empty line
			j := i + 1
			for j < len(lines) && strings.TrimSpace(lines[j]) != "" {
				block = append(block, lines[j])
				j++
			}

			// Add the block to each hashtag
			blockText := strings.Join(block, "\n")
			for _, hashtag := range hashtags {
				result[hashtag] = append(result[hashtag], blockText)
			}

			// Skip the lines we've already processed
			i = j - 1
		}
	}

	return result
}
