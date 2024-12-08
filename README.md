# sigolyze
[![Go](https://github.com/imyazip/sigolyze/actions/workflows/tests.yml/badge.svg)](https://github.com/imyazip/sigolyze/actions/workflows/tests.yml)

`sigolyze` is a Go package designed for efficient signature-based substring searching using the Aho-Corasick algorithm and regular expression matching. This package allows you to quickly search for predefined patterns in data, both as plain text and regular expressions. It supports using multiple tags to categorize patterns, making it easier to filter and group them.

## Features

- **Signature-based searching**: Utilizes the Aho-Corasick algorithm for fast matching of exact string patterns.
- **Regex support**: Supports regular expressions alongside plain string patterns for flexible searches.
- **Tag-based categorization**: Allows patterns to be categorized with tags for easier grouping and filtering.
- **Efficient matching**: Can match multiple signatures in a single search using both Aho-Corasick and regular expressions.
- **JSON-based signature loading**: Signatures can be defined in JSON format and loaded into the compiler for searching.

## Installation

To install the `sigolyze` package, use the Go module system:

```bash
go get github.com/yourusername/sigolyze
```

## Usage

### Compiler setup
You can create a new compiler instance with NewCompiler, which will hold all your compiled signatures.
```go
compiler := sigolyze.NewCompiler()
```

### Loading Signatures
Signatures can be loaded from a JSON byte slice or from a .json file. Each signature consists of patterns that can be either plain strings or regular expressions.

#### Load Signature from JSON data
```go
signatureData := []byte(`{
    "name": "example_signature",
    "patterns": [{"name": "pattern1", "value": "abc", "is_regex": false}],
    "tags": ["tag1", "tag2"],
    "meta": [{"name": "info", "info": ["some info"]}]
}`)
err := compiler.LoadSignature(signatureData)
if err != nil {
    fmt.Println("Error loading signature:", err)
}
```

#### Load Signature from a `.json` file
```go
err := compiler.LoadSignatureFromJson("path/to/signature.json")
if err != nil {
    fmt.Println("Error loading signature from file:", err)
}
```

### Searching for Matches
Once signatures are loaded, you can perform searches on a given data string. The Match function will return all signatures that match the data.

```go
data := "sample data containing abc"
matches := sigolyze.Match(compiler, data)

for _, match := range matches {
    fmt.Println("Matched Signature:", match.Name)
}
```
If you want to search for matches based on specific tags, use the MatchTags function:

```go
tags := []string{"tag1"}
matches := sigolyze.MatchTags(compiler, data, tags)

for _, match := range matches {
    fmt.Println("Matched Signature with tag:", match.Name)
}
```

### Signature structure

Each signature consists of the following components:

- `Name`: The name of the signature.
- `Patterns`: A list of patterns to search for, each pattern having
  - `Name`: The name of the pattern.
  - `Value`: The value of the pattern (can be a string or regex).
  - `IsRegex`: Boolean indicating whether the pattern is a regular expression.
- `Tags`: A list of tags used to categorize the signature.
- `Meta`: Additional meta information for the signature.
- `Matcher`: Precompiled matcher for efficient search.

### Example `.json`
```json
{
    "name": "Example Signature",
    "patterns": [
      {
        "name": "Pattern 1",
        "value": "Value[1-9]+",
        "is_regex": true
      },
      {
        "name": "Pattern 2",
        "value": "Value2",
        "is_regex": false
      }
    ],
    "tags": ["tag1", "tag2"],
    "meta": [
      {
        "name": "Meta1",
        "info": ["detail1", "detail2"]
      },
      {
        "name": "Meta2",
        "info": ["detail3"]
      }
    ]
  }

```