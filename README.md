# Markdownsift

Markdownsift extracts tagged blocks of text from markdown files and writes them order by date to stdout or files.

## Usage

```
mds --help
```

## How it works

### Paragraphs

A file `2025-07-20.md` that contains a paragraph tagged with `#tougths` 

```
Started the year with ambitious plans. #thoughts
```
will be extracted into a topic file `thougts.md`:

```
# Content tagged with #thoughts

2025-07-20:
Started the year with ambitious plans. #thoughts
```
A tagged paragraph ends with a newline, thus everthing before the tagged line and everthing after the new line is not exctracted.

## TODOso

- [ ] Make file pattern configurable
- [ ] Sort blocks in descent order, given by cli flag
