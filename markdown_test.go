package markdownsift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCollectHashtaggedContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string][]Block
	}{
		{
			name: "single hashtag with block",
			input: `Some text without hashtags
This line has #work
This belongs to work block
And this too

This line is separate`,
			expected: map[string][]Block{
				"#work": {{Content: "This line has #work\nThis belongs to work block\nAnd this too"}},
			},
		},
		{
			name: "multiple hashtags in same line",
			input: `Meeting about #work and #project
Notes from the meeting
Action items

Another line`,
			expected: map[string][]Block{
				"#work":    {{Content: "Meeting about #work and #project\nNotes from the meeting\nAction items"}},
				"#project": {{Content: "Meeting about #work and #project\nNotes from the meeting\nAction items"}},
			},
		},
		{
			name: "multiple blocks for same hashtag",
			input: `First #work block
Content of first block

Second #work block
Content of second block

No hashtag here`,
			expected: map[string][]Block{
				"#work": {
					{Content: "First #work block\nContent of first block"},
					{Content: "Second #work block\nContent of second block"},
				},
			},
		},
		{
			name: "hashtag at end of file without trailing newline",
			input: `Last #task without newline
Final content`,
			expected: map[string][]Block{
				"#task": {{Content: "Last #task without newline\nFinal content"}},
			},
		},
		{
			name: "single line hashtag blocks",
			input: `Single #quick note

Another #fast item

#urgent task`,
			expected: map[string][]Block{
				"#quick":  {{Content: "Single #quick note"}},
				"#fast":   {{Content: "Another #fast item"}},
				"#urgent": {{Content: "#urgent task"}},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: map[string][]Block{},
		},
		{
			name: "no hashtags",
			input: `Just regular text
No tags here
Multiple lines`,
			expected: map[string][]Block{},
		},
		{
			name: "hashtag with special characters",
			input: `Task #work_project with underscore
Details about the task

Meeting #meeting123 with numbers
Meeting notes`,
			expected: map[string][]Block{
				"#work_project": {{Content: "Task #work_project with underscore\nDetails about the task"}},
				"#meeting123":   {{Content: "Meeting #meeting123 with numbers\nMeeting notes"}},
			},
		},
		{
			name: "empty lines between blocks",
			input: `First #block here
Block content


Second #block there
More content



Third #block elsewhere
Final content`,
			expected: map[string][]Block{
				"#block": {
					{Content: "First #block here\nBlock content"},
					{Content: "Second #block there\nMore content"},
					{Content: "Third #block elsewhere\nFinal content"},
				},
			},
		},
		{
			name: "headed blocks with empty lines preserved",
			input: `# section #work
Content for work

More work content

Another line

# section #personal
Personal content

With empty lines`,
			expected: map[string][]Block{
				"#work":     {{Content: "# section #work\nContent for work\n\nMore work content\n\nAnother line\n"}},
				"#personal": {{Content: "# section #personal\nPersonal content\n\nWith empty lines"}},
			},
		},
		{
			name: "mixed regular and headed blocks",
			input: `Regular #task item
Task details

# heading #project
Project description

More project info

Another #task item
More task details

# notes #meeting
Meeting content

Final notes`,
			expected: map[string][]Block{
				"#task": {
					{Content: "Regular #task item\nTask details"},
				},
				"#project": {{Content: "# heading #project\nProject description\n\nMore project info\n\nAnother #task item\nMore task details\n"}},
				"#meeting": {{Content: "# notes #meeting\nMeeting content\n\nFinal notes"}},
			},
		},
		{
			name: "headed block at document end",
			input: `# section #final
Last content

With empty line`,
			expected: map[string][]Block{
				"#final": {{Content: "# section #final\nLast content\n\nWith empty line"}},
			},
		},
		{
			name: "multiple hashtags in headed block",
			input: `# priority #work #urgent
Important task content

More details

# notes #personal
Personal content`,
			expected: map[string][]Block{
				"#work":     {{Content: "# priority #work #urgent\nImportant task content\n\nMore details\n"}},
				"#urgent":   {{Content: "# priority #work #urgent\nImportant task content\n\nMore details\n"}},
				"#personal": {{Content: "# notes #personal\nPersonal content"}},
			},
		},
		{
			name: "markdown headings as headed blocks",
			input: `Regular #task item
Task details

## Important Section #ai
Content under the heading

More content with empty lines

Some final content

### Another Section #work
Final section content`,
			expected: map[string][]Block{
				"#task": {{Content: "Regular #task item\nTask details"}},
				"#ai":   {{Content: "## Important Section #ai\nContent under the heading\n\nMore content with empty lines\n\nSome final content\n"}},
				"#work": {{Content: "### Another Section #work\nFinal section content"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result := collectHashtaggedContent(reader)

			// Check if we have the expected number of hashtags
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d hashtags, got %d", len(tt.expected), len(result))
				return
			}

			// Check each hashtag and its blocks
			for hashtag, expectedBlocks := range tt.expected {
				actualBlocks, exists := result[hashtag]
				if !exists {
					t.Errorf("Expected hashtag %s not found in result", hashtag)
					continue
				}

				if len(actualBlocks) != len(expectedBlocks) {
					t.Errorf("For hashtag %s: expected %d blocks, got %d",
						hashtag, len(expectedBlocks), len(actualBlocks))
					continue
				}

				for i, expectedBlock := range expectedBlocks {
					if actualBlocks[i].Content != expectedBlock.Content {
						t.Errorf("For hashtag %s, block %d:\nExpected:\n%s\nGot:\n%s",
							hashtag, i, expectedBlock.Content, actualBlocks[i].Content)
					}
				}
			}
		})
	}
}

func TestExtractWithReader(t *testing.T) {
	// Test with different types of readers to ensure io.Reader interface works
	input := `Test #example block
Content here

Another #test item
More content`

	reader := strings.NewReader(input)
	result := collectHashtaggedContent(reader)

	if len(result) != 2 {
		t.Errorf("Expected 2 hashtags, got %d", len(result))
	}

	if _, exists := result["#example"]; !exists {
		t.Error("Expected #example hashtag not found")
	}

	if _, exists := result["#test"]; !exists {
		t.Error("Expected #test hashtag not found")
	}
}

func TestWriteSnippets(t *testing.T) {
	// Sample test data
	snippets := map[string][]Block{
		"#work": {
			{Content: "Task 1 content\nMore details", Date: time.Now()},
			{Content: "Task 2 content", Date: time.Now()},
		},
		"#personal": {
			{Content: "Personal note\nWith multiple lines", Date: time.Now()},
		},
	}

	t.Run("write to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		WriteSnippets(&buf, snippets, nil)

		output := buf.String()

		// Check that both hashtags are present
		if !strings.Contains(output, "#work:") {
			t.Error("Expected #work hashtag in output")
		}
		if !strings.Contains(output, "#personal:") {
			t.Error("Expected #personal hashtag in output")
		}

		// Check that block numbers are present
		if !strings.Contains(output, "Block 1 (") {
			t.Error("Expected Block 1 in output")
		}
		if !strings.Contains(output, "Block 2 (") {
			t.Error("Expected Block 2 in output")
		}

		// Check that content is present
		if !strings.Contains(output, "Task 1 content") {
			t.Error("Expected task content in output")
		}
		if !strings.Contains(output, "Personal note") {
			t.Error("Expected personal note in output")
		}
	})

	t.Run("write to string builder", func(t *testing.T) {
		var builder strings.Builder
		WriteSnippets(&builder, snippets, nil)

		output := builder.String()

		// Verify the output contains expected structure
		lines := strings.Split(output, "\n")
		var hashtagLines []string
		for _, line := range lines {
			if strings.HasSuffix(line, ":") && strings.HasPrefix(line, "#") {
				hashtagLines = append(hashtagLines, line)
			}
		}

		// Should have 2 hashtag headers
		if len(hashtagLines) != 2 {
			t.Errorf("Expected 2 hashtag headers, got %d", len(hashtagLines))
		}
	})

	t.Run("empty snippets", func(t *testing.T) {
		var buf bytes.Buffer
		emptySnippets := make(map[string][]Block)
		WriteSnippets(&buf, emptySnippets, nil)

		if buf.Len() != 0 {
			t.Error("Expected empty output for empty snippets")
		}
	})

	t.Run("single hashtag single block", func(t *testing.T) {
		var buf bytes.Buffer
		singleSnippet := map[string][]Block{
			"#test": {{Content: "Single line content", Date: time.Now()}},
		}
		WriteSnippets(&buf, singleSnippet, nil)

		output := buf.String()
		// Check that the output contains expected elements rather than exact match
		// since date format will vary
		if !strings.Contains(output, "#test:") {
			t.Error("Expected #test hashtag in output")
		}
		if !strings.Contains(output, "Block 1 (") {
			t.Error("Expected 'Block 1 (' with date in output")
		}
		if !strings.Contains(output, "Single line content") {
			t.Error("Expected content in output")
		}
	})

	t.Run("filter by specific tags", func(t *testing.T) {
		var buf bytes.Buffer
		testSnippets := map[string][]Block{
			"#work": {
				{Content: "Work task 1", Date: time.Now()},
				{Content: "Work task 2", Date: time.Now()},
			},
			"#personal": {
				{Content: "Personal note 1", Date: time.Now()},
			},
			"#project": {
				{Content: "Project update", Date: time.Now()},
			},
		}

		// Only write #work and #project tags
		tags := []string{"#work", "#project"}
		WriteSnippets(&buf, testSnippets, tags)

		output := buf.String()

		// Should contain #work and #project
		if !strings.Contains(output, "#work:") {
			t.Error("Expected #work hashtag in filtered output")
		}
		if !strings.Contains(output, "#project:") {
			t.Error("Expected #project hashtag in filtered output")
		}
		if !strings.Contains(output, "Work task 1") {
			t.Error("Expected work content in filtered output")
		}
		if !strings.Contains(output, "Project update") {
			t.Error("Expected project content in filtered output")
		}

		// Should NOT contain #personal
		if strings.Contains(output, "#personal:") {
			t.Error("Should not contain #personal hashtag in filtered output")
		}
		if strings.Contains(output, "Personal note 1") {
			t.Error("Should not contain personal content in filtered output")
		}
	})
}
