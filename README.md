# go-highlight

## Description

A somewhat crude syntax highlighter loosely based on [Pygments](pygments.org).

*See [chroma](https://github.com/alecthomas/chroma) for a newer, better, and far more complete syntax highlighter!*

## Installation

go-highlight can be installed with the regular `go get` command:

    go get github.com/johnsto/go-highlight

Lexers are stored in a separate package, so remember to install those, too:

    go get github.com/johnsto/go-highlight/lexers

The CLI app can also be installed with `go get`:

    go get github.com/johnsto/go-highlight/cmd/highlight

Run tests:

    go test github.com/johnsto/go-highlight

## Usage

Importing for use in your code is as simple as importing both the base
`highlight` package, and registering the default lexers as an anonymous
import:

```go
import "github.com/johnsto/go-highlight"
import _ "github.com/johnsto/go-highlight/lexers"
```

Tokenizers can be retrieved by content type or filename:

```go
tokenizer, err = highlight.GetTokenizerForContentType("application/json")
// or:
tokenizer, err = highlight.GetTokenizerForFilename("futurama.json")
```

Use `Tokenize` or `TokenizeString` to tokenize an `io.Reader` or `string`
respectively:

```go
err = tokenizer.Tokenize(reader, func(t highlight.Token) error {
	_, err := fmt.Printf(t.Value)
	return err
})
```

For colourised terminal output, import the "term" package:

```go
import "github.com/johnsto/go-highlight/output/term"
```

Then instantiate the emitter and tokenize to it:

```go
emitter = output.NewDebugOutputter()
err = tokenizer.Tokenize(reader, emitter.Emit)
```

Or for formatted (indented) output, use `Format` instead:

```go
err = tokenizer.Format(reader, emitter.Emit)
```
