# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library called `markdownsift` that extracts hashtag-grouped content blocks from markdown files. The main functionality is provided by the `CollectSnippets` function which processes multiple markdown files in a directory and groups content blocks by their hashtags.

## Architecture

- **Core library**: `markdown.go` contains the main `CollectSnippets` function and internal parsing logic
- **CLI application**: `cmd/markdownsift/main.go` provides a command-line interface
- **Testing**: `markdown_test.go` contains comprehensive test cases

The library works by:
1. `CollectSnippets` scans a directory for markdown files matching the default pattern
2. For each file, `collectHashtaggedContent` reads content and finds hashtags (using regex `#\w+`)
3. Content blocks following hashtag lines are collected until an empty line
4. Returns a consolidated map where keys are hashtags and values are arrays of content blocks from all files

## Common Commands

### Build
```bash
go build ./cmd/markdownsift
```

### Run Tests
```bash
go test
```

### Run Single Test
```bash
go test -run TestExtract
```

### Run CLI
```bash
# Use default path
go run ./cmd/markdownsift/main.go

# Specify custom path
go run ./cmd/markdownsift/main.go -path /path/to/your/markdown/files
```

## Development Notes

- The regex pattern `#\w+` matches hashtags with word characters (letters, digits, underscore)
- Content blocks are collected until an empty line is encountered
- Multiple hashtags on the same line will add the entire block to each hashtag's collection

## Workflow

- always review my last commit and give me your recommendations