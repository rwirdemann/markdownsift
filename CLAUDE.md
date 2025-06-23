# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library called `markdownsift` that extracts hashtag-grouped content blocks from markdown files. The main functionality is provided by the `CollectHashtaggedContent` function which parses markdown text and groups content blocks by their hashtags.

## Architecture

- **Core library**: `markdown.go` contains the main `CollectHashtaggedContent` function
- **CLI application**: `cmd/markdownsift/main.go` provides a command-line interface
- **Testing**: `markdown_test.go` contains comprehensive test cases

The library works by:
1. Reading markdown content from an `io.Reader`
2. Finding lines containing hashtags (using regex `#\w+`)
3. Collecting content blocks that follow hashtag lines until an empty line
4. Returning a map where keys are hashtags and values are arrays of content blocks

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
go run ./cmd/markdownsift/main.go
```

## Development Notes

- The main function currently has a hardcoded file path that should be made configurable
- The regex pattern `#\w+` matches hashtags with word characters (letters, digits, underscore)
- Content blocks are collected until an empty line is encountered
- Multiple hashtags on the same line will add the entire block to each hashtag's collection