# go-regex
A simple command-line tool for testing regular expressions in Go.

## Installation

You can install the tool with `go install`:

```bash
go install github.com/yourusername/go-regex/cmd/go_regex_cli@latest
```

## Usage
```bash
go-regex -r "your_regex" -s "string_to_test"
```
This will return `true` if the regex matches the string, and `false` otherwise

## Example
```bash
go-regex -r "a*b" -s "aaaab"
true
```
