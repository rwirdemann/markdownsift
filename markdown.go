package markdownsift

import (
	"io"
	"regexp"
	"strings"
)

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
