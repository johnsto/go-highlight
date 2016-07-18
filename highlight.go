package highlight

import (
	"fmt"
	"io"
)

var Debug bool = false

// Token represents one item of parsed output, containing the parsed value
// and its detected type.
type Token struct {
	Value string
	Type  TokenType
	State string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s:%#v {%s})", t.Type, t.Value, t.State)
}

// Tokenizer represents a type capable of tokenizing data from an input
// source.
type Tokenizer interface {
	// Tokenize reads from the given input and emits tokens to the output
	// channel. Will end on any error from the reader, including io.EOF to
	// signify the end of input.
	Tokenize(io.Reader, func(Token) error) error

	// AcceptsFilename returns true if this Lexer thinks it is suitable for
	// the given filename. An error will be returned iff an invalid filename
	// pattern is registered by the Lexer.
	AcceptsFilename(name string) (bool, error)

	// AcceptsMediaType returns true if this Lexer thinks it is suitable for
	// the given meda (MIME) type. An error wil be returned iff the given mime
	// type is invalid.
	AcceptsMediaType(name string) (bool, error)
}
