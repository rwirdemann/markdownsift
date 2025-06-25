package markdownsift

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// DefaultPattern is the default pattern used to match markdown files.
	DefaultPattern = "^\\d{4}-\\d{2}-\\d{2}\\.md$"
)

// CollectSnippets scans the specified directory and returns a map of hashtags to associated content snippets.
// The hashtags and content are extracted from the files that match the predefined pattern in the directory.
func CollectSnippets(path string) map[string][]string {
	files, err := listFiles(path)
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		return nil
	}

	var snippets = map[string][]string{}
	for _, file := range files {
		func() {
			fmt.Printf("Processing file: %s\n", file)
			file, err := os.Open(file)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			result := collectHashtaggedContent(file)
			for tag, blocks := range result {
				snippets[tag] = append(snippets[tag], blocks...)
			}
		}()
	}
	return snippets
}

// listFiles returns a list of files in the given path that match the DefaultPattern.
func listFiles(path string) ([]string, error) {

	var matchingFiles []string

	// Compile the regex pattern
	regex, err := regexp.Compile(DefaultPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern '%s': %w", DefaultPattern, err)
	}

	// Read the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", path, err)
	}

	// Filter files that match the pattern
	for _, entry := range entries {
		if !entry.IsDir() && regex.MatchString(entry.Name()) {
			matchingFiles = append(matchingFiles, filepath.Join(path, entry.Name()))
		}
	}

	return matchingFiles, nil
}

// collectHashtaggedContent reads content from the given reader and returns a map of tags pointing to snippets tagged
// with the hashtag.
func collectHashtaggedContent(reader io.Reader) map[string][]string {
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
