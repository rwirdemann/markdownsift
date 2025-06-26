# markdownsift

A Go library for extracting and organizing hashtag-grouped content blocks from markdown files.

## Overview

`markdownsift` scans markdown files in a directory and extracts content blocks that are associated with hashtags (like `#work`, `#personal`, `#urgent`). Each content block is now stored with its content and the date it was processed, making it easy to track when information was collected.

## Features

- **Hashtag-based content extraction**: Finds and groups content by hashtags
- **Timestamped blocks**: Each content block includes the date/time it was processed
- **Two block types supported**:
  - **Regular blocks**: Content following a hashtag until an empty line
  - **Headed blocks**: Markdown sections starting with headings (`# ## ### ####`) that contain hashtags
- **Command-line interface**: Easy-to-use CLI for processing markdown files
- **Tag filtering**: Process only specific hashtags you're interested in
- **File pattern matching**: Automatically processes files matching `YYYY-MM-DD.md` format

## Installation

```bash
go get github.com/rwirdemann/markdownsift
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "os"
    "github.com/rwirdemann/markdownsift"
)

func main() {
    // Collect all hashtag blocks from markdown files
    snippets := markdownsift.CollectSnippets("path/to/markdown/files")
    
    // Write all blocks to stdout
    markdownsift.WriteSnippets(os.Stdout, snippets, nil)
    
    // Or filter by specific tags
    tags := []string{"#work", "#urgent"}
    markdownsift.WriteSnippets(os.Stdout, snippets, tags)
}
```

### Block Structure

Each content block is now represented by a `Block` struct:

```go
type Block struct {
    Content string    // The actual content text
    Date    time.Time // When the block was processed (set to current time)
}
```

### Command Line Interface

Build the CLI tool:

```bash
go build ./cmd/markdownsift
```

Basic usage:

```bash
# Process all markdown files in default directory
./markdownsift

# Specify custom directory
./markdownsift -path /path/to/your/markdown/files

# Filter by specific tags (omit # symbol)
./markdownsift -path /path/to/files -tags work,urgent,personal
```

### Example Input

```markdown
# Daily Notes - 2025-06-26

## Morning Tasks #work
Review pull requests
Update documentation
Team standup at 9:30

## Personal Reminder #personal
Buy groceries after work
Call dentist to schedule appointment

## Critical Issue #urgent #work  
Server downtime reported
Need immediate investigation
```

### Example Output

```
#work:
Block 1 (2025-06-26 14:30:15):
## Morning Tasks #work
Review pull requests
Update documentation
Team standup at 9:30

Block 2 (2025-06-26 14:30:15):
## Critical Issue #urgent #work  
Server downtime reported
Need immediate investigation

#personal:
Block 1 (2025-06-26 14:30:15):
## Personal Reminder #personal
Buy groceries after work
Call dentist to schedule appointment

#urgent:
Block 1 (2025-06-26 14:30:15):
## Critical Issue #urgent #work  
Server downtime reported
Need immediate investigation
```

## How It Works

1. **File Discovery**: Scans the specified directory for files matching the pattern `YYYY-MM-DD.md`
2. **Content Parsing**: For each file, finds lines containing hashtags using regex `#\w+`
3. **Block Collection**: 
   - For regular blocks: Collects content until an empty line
   - For headed blocks: Collects content until the next heading or end of file
4. **Timestamping**: Each block is stamped with the current date/time when processed
5. **Grouping**: Organizes blocks by hashtag in a map structure

## Development

### Running Tests

```bash
go test
```

### Running a Specific Test

```bash
go test -run TestCollectHashtaggedContent
```

### Building

```bash
go build ./cmd/markdownsift
```

## API Reference

### Functions

#### `CollectSnippets(path string) map[string][]Block`

Scans the specified directory and returns a map of hashtags to associated content blocks.

#### `WriteSnippets(writer io.Writer, snippets map[string][]Block, tags []string)`

Writes the snippets matching the given tags to the specified writer. If `tags` is empty or nil, all snippets are written.

### Types

#### `Block`

```go
type Block struct {
    Content string    // The content text of the block
    Date    time.Time // Timestamp when the block was processed
}
```

## File Pattern

By default, the library processes files matching the pattern `^\\d{4}-\\d{2}-\\d{2}\\.md$` (e.g., `2025-06-26.md`). This pattern is designed for daily note-taking systems but can be customized by modifying the `DefaultPattern` constant.

## License

This project is open source. Please check the license file for details.