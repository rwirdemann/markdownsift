package markdownsift

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	// DefaultPattern is the default pattern used to match markdown files.
	DefaultPattern = "^\\d{4}-\\d{2}-\\d{2}\\.md$"
)

// Block represents a content block with its associated date
type Block struct {
	Content string
	Date    time.Time
}

// CollectSnippets scans the specified directory and returns a map of hashtags to associated content snippets. The
// hashtags and content are extracted from the files that match the predefined pattern in the directory.
func CollectSnippets(path string) map[string][]Block {
	files, err := listFiles(path)
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		return nil
	}

	var snippets = map[string][]Block{}
	for _, file := range files {
		func() {
			// Parse date from filename (format: YYYY-MM-DD.md)
			filename := filepath.Base(file)
			dateStr := strings.TrimSuffix(filename, ".md")
			fileDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				// Fallback to current time if parsing fails
				fileDate = time.Now()
			}

			f, err := os.Open(file)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			result := collectHashtaggedContent(f, fileDate)
			for tag, blocks := range result {
				snippets[tag] = append(snippets[tag], blocks...)
			}
		}()
	}
	return snippets
}

// WriteSnippets writes the snippets matching the given tags to the specified writer.
func WriteSnippets(writer io.Writer, snippets map[string][]Block, tags []string) {
	for tag, blocks := range snippets {
		fmt.Fprintf(writer, "# Content tagged by %s\n", tag)
		for _, block := range blocks {
			fmt.Fprintf(writer, "%s:\n%s\n\n", block.Date.Format("2006-01-02"), block.Content)
		}
	}
}

func Filter(snippets map[string][]Block, tags []string) map[string][]Block {
	if len(tags) == 0 {
		return snippets
	}

	var filtered = make(map[string][]Block)
	for tag, blocks := range snippets {
		if slices.Contains(tags, tag) {
			filtered[tag] = blocks
		}
	}
	return filtered
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
// with the hashtag. It handles both regular blocks (ending at empty lines) and headed blocks (ending at next headed
// block or end of document).
func collectHashtaggedContent(reader io.Reader, date time.Time) map[string][]Block {
	result := make(map[string][]Block)

	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return result
	}

	lines := strings.Split(string(content), "\n")
	hashtagRegex := regexp.MustCompile(`#\w+`)

	// First pass: identify headed blocks (hashtag lines that start document sections)
	headedBlockStarts := make(map[int]bool)
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if len(line) > 0 && hashtagRegex.MatchString(line) {
			// Check if this looks like a markdown heading (starts with # ## ### ####)
			if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") || strings.HasPrefix(line, "#### ") {
				headedBlockStarts[i] = true
			}
		}
	}

	// Second pass: collect blocks with appropriate termination logic
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		hashtags := hashtagRegex.FindAllString(line, -1)

		if len(hashtags) > 0 {
			block := []string{line}
			j := i + 1

			if headedBlockStarts[i] {
				// This is a headed block - collect until next headed block or end of document
				for j < len(lines) {
					if headedBlockStarts[j] {
						break // Stop at next headed block
					}
					block = append(block, lines[j])
					j++
				}
			} else {
				// This is a regular block - collect until empty line
				for j < len(lines) && strings.TrimSpace(lines[j]) != "" {
					block = append(block, lines[j])
					j++
				}
			}

			// Add the block to each hashtag
			blockText := strings.Join(block, "\n")
			blockInstance := Block{
				Content: blockText,
				Date:    date,
			}
			for _, hashtag := range hashtags {
				result[hashtag] = append(result[hashtag], blockInstance)
			}

			// Skip the lines we've already processed
			i = j - 1
		}
	}

	return result
}
