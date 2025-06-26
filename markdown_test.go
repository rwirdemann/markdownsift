package markdownsift

import (
	"strings"
	"testing"
)

func TestCollectHashtaggedContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string][]string
	}{
		{
			name: "single hashtag with block",
			input: `Some text without hashtags
This line has #work
This belongs to work block
And this too

This line is separate`,
			expected: map[string][]string{
				"#work": {"This line has #work\nThis belongs to work block\nAnd this too"},
			},
		},
		{
			name: "multiple hashtags in same line",
			input: `Meeting about #work and #project
Notes from the meeting
Action items

Another line`,
			expected: map[string][]string{
				"#work":    {"Meeting about #work and #project\nNotes from the meeting\nAction items"},
				"#project": {"Meeting about #work and #project\nNotes from the meeting\nAction items"},
			},
		},
		{
			name: "multiple blocks for same hashtag",
			input: `First #work block
Content of first block

Second #work block
Content of second block

No hashtag here`,
			expected: map[string][]string{
				"#work": {
					"First #work block\nContent of first block",
					"Second #work block\nContent of second block",
				},
			},
		},
		{
			name: "hashtag at end of file without trailing newline",
			input: `Last #task without newline
Final content`,
			expected: map[string][]string{
				"#task": {"Last #task without newline\nFinal content"},
			},
		},
		{
			name: "single line hashtag blocks",
			input: `Single #quick note

Another #fast item

#urgent task`,
			expected: map[string][]string{
				"#quick":  {"Single #quick note"},
				"#fast":   {"Another #fast item"},
				"#urgent": {"#urgent task"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: map[string][]string{},
		},
		{
			name: "no hashtags",
			input: `Just regular text
No tags here
Multiple lines`,
			expected: map[string][]string{},
		},
		{
			name: "hashtag with special characters",
			input: `Task #work_project with underscore
Details about the task

Meeting #meeting123 with numbers
Meeting notes`,
			expected: map[string][]string{
				"#work_project": {"Task #work_project with underscore\nDetails about the task"},
				"#meeting123":   {"Meeting #meeting123 with numbers\nMeeting notes"},
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
			expected: map[string][]string{
				"#block": {
					"First #block here\nBlock content",
					"Second #block there\nMore content",
					"Third #block elsewhere\nFinal content",
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
			expected: map[string][]string{
				"#work":     {"# section #work\nContent for work\n\nMore work content\n\nAnother line\n"},
				"#personal": {"# section #personal\nPersonal content\n\nWith empty lines"},
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
			expected: map[string][]string{
				"#task": {
					"Regular #task item\nTask details",
				},
				"#project": {"# heading #project\nProject description\n\nMore project info\n\nAnother #task item\nMore task details\n"},
				"#meeting": {"# notes #meeting\nMeeting content\n\nFinal notes"},
			},
		},
		{
			name: "headed block at document end",
			input: `# section #final
Last content

With empty line`,
			expected: map[string][]string{
				"#final": {"# section #final\nLast content\n\nWith empty line"},
			},
		},
		{
			name: "multiple hashtags in headed block",
			input: `# priority #work #urgent
Important task content

More details

# notes #personal
Personal content`,
			expected: map[string][]string{
				"#work":     {"# priority #work #urgent\nImportant task content\n\nMore details\n"},
				"#urgent":   {"# priority #work #urgent\nImportant task content\n\nMore details\n"},
				"#personal": {"# notes #personal\nPersonal content"},
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
			expected: map[string][]string{
				"#task": {"Regular #task item\nTask details"},
				"#ai":   {"## Important Section #ai\nContent under the heading\n\nMore content with empty lines\n\nSome final content\n"},
				"#work": {"### Another Section #work\nFinal section content"},
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
					if actualBlocks[i] != expectedBlock {
						t.Errorf("For hashtag %s, block %d:\nExpected:\n%s\nGot:\n%s",
							hashtag, i, expectedBlock, actualBlocks[i])
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
