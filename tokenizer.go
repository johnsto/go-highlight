package highlight

import "bufio"

// Tokenizer represents a type capable of tokenizing data from an input
// source.
type Tokenizer interface {
	// Tokenize reads from the given input and emits tokens to the output
	// channel. Will end on any error from the reader, including io.EOF to
	// signify the end of input.
	Tokenize(*bufio.Reader, func(Token) error) error

	// Format behaves exactly as Tokenize, except it also formats the output.
	Format(*bufio.Reader, func(Token) error) error

	// AcceptsFilename returns true if this Lexer thinks it is suitable for
	// the given filename. An error will be returned iff an invalid filename
	// pattern is registered by the Lexer.
	AcceptsFilename(name string) (bool, error)

	// AcceptsMediaType returns true if this Lexer thinks it is suitable for
	// the given meda (MIME) type. An error wil be returned iff the given mime
	// type is invalid.
	AcceptsMediaType(name string) (bool, error)

	// ListMediaTypes lists the media types this Tokenizer advertises support
	// for, e.g. ["application/json"]
	ListMediaTypes() []string

	// ListFilenames lists the filename patterns this Tokenizer advertises
	// support for, e.g. ["*.json"]
	ListFilenames() []string
}
