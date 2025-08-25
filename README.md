# Markdownsift

Markdownsift extracts tagged blocks of text from markdown files and writes them order by date to stdout or files.

## Usage

```
mds --help
```

## How it works

### Paragraphs

A file `2025-07-20.md` that contains a paragraph tagged with `#thoughts` 

```
Started the year with ambitious plans. #thoughts
```
will be extracted into a topic file `thoughts.md`:

```
# Content tagged with #thoughts

2025-07-20:
Started the year with ambitious plans. #thoughts
```
A tagged paragraph ends with a newline, thus everything before the tagged line and everything after the new line is not extracted.

## TODOs

- [ ] Make file pattern configurable
- [ ] Sort blocks in descent order, given by cli flag
