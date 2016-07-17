package highlight

import "io"

// Filter describes a type that is capable of filtering/processing tokens.
type Filter interface {
	// Filter reads tokens from `in` and outputs tokens to `out`, typically
	// modifying or filtering them along the way. The function should return
	// as soon as the input is exhausted (i.e. the channel is closed), or an
	// error is encountered.
	Filter(lexer Lexer, out func(Token) error) func(Token) error
}

// FilterFunc is a helper type allowing filter functions to be used as
// filters.
type FilterFunc func(lexer Lexer, out func(Token) error) func(Token) error

func (f FilterFunc) Filter(lexer Lexer, out func(Token) error) func(Token) error {
	return f(lexer, out)
}

type Filters []Filter

// Filter runs the input through each filter in series, emitting the final
// result to `out`. This function will return as soon as the last token has
// been processed, or iff an error is encountered by one of the filters.
//
// It is safe to close the output channel as soon as this function returns.
func (fs Filters) Filter(lexer Lexer, out func(Token) error) func(Token) error {
	for _, f := range fs {
		out = f.Filter(lexer, out)
	}
	return out
}

// PassthroughFilter simply emits each token to the output without
// modification.
var PassthroughFilter = FilterFunc(
	func(l Lexer, out func(Token) error) func(Token) error {
		return func(t Token) error {
			return out(t)
		}
	})

// RemoveEmptiesFilter removes empty (zero-length) tokens from the output.
var RemoveEmptiesFilter = FilterFunc(
	func(l Lexer, out func(Token) error) func(Token) error {
		return func(t Token) error {
			if t.Value != "" {
				return out(t)
			}
			return nil
		}
	})

// MergeTokensFilter combines Tokens if they have the same type.
var MergeTokensFilter = FilterFunc(
	func(lexer Lexer, out func(Token) error) func(Token) error {
		curr := Token{}

		return func(t Token) error {
			if t.Type == "" {
				out(curr)
				return io.EOF
			} else if t.Type == curr.Type {
				// Same as last token; combine
				curr.Value += t.Value
				return nil
			} else if curr.Value != "" {
				out(curr)
			}
			curr = Token{
				Value: t.Value,
				Type:  t.Type,
				State: t.State,
			}
			return nil
		}
	})
