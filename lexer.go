package highlight

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"path"
	"strings"
)

// Lexer defines a simple state-based lexer.
type Lexer struct {
	Name      string
	States    States
	Filters   Filters
	Formatter Filter
	Filenames []string
	MimeTypes []string
}

func (l Lexer) Format(r *bufio.Reader, emit func(Token) error) error {
	if l.Formatter == nil {
		return l.Tokenize(r, emit)
	}
	return l.Tokenize(r, l.Formatter.Filter(emit))
}

// Tokenize reads from the given input and emits tokens to the output channel.
// Will end on any error from the reader, including io.EOF to signify the end
// of input.
func (l Lexer) Tokenize(br *bufio.Reader, emit func(Token) error) error {
	states, err := l.States.Compile()
	if err != nil {
		return err
	}

	emit = l.Filters.Filter(emit)

	stack := &Stack{"root"}
	eol := false
	var subject = ""
	for {
		next, err := br.ReadString('\n')

		if err == bufio.ErrBufferFull {
			eol = false
		} else if err == io.EOF {
			eol = true
		} else if err != nil {
			return emit(EndToken)
		} else {
			eol = strings.HasSuffix(next, "\n")
		}

		subject = subject + next

		if subject == "" && err == io.EOF {
			emit(EndToken)
			return err
		}

		for subject != "" {
			// Match current state against current subject
			stateName := stack.Peek()
			state := states.Get(stateName)

			// Tokenize input
			n, rule, tokens, err := state.Match(subject)
			if err != nil {
				return emit(EndToken)
			}

			// No rules matched
			if rule == nil {
				if !eol {
					// Read more data for the current line
					break
				} else {
					// Emit entire subject as an error
					tokens = []Token{{Value: subject, Type: Error}}
					n = len(subject)
				}
			}

			// Emit each token to the output
			for _, t := range tokens {
				t.State = stateName
				if err := emit(t); err != nil {
					emit(EndToken)
					return err
				}
			}

			// Update state
			if rule == nil {
				if !eol {
					// Didn't match at all, reset to root state
					stack.Empty()
					stack.Push("root")
				}
			} else {
				// Push new states as appropriate
				for _, state := range rule.Stack() {
					if state == "#pop" {
						stack.Pop()
					} else if state != "" {
						stack.Push(state)
					}
				}
			}

			if stack.Len() == 0 {
				return emit(EndToken)
			}

			// Consume matched part
			subject = subject[n:]
		}
	}

	return nil
}

// TokenizeString is a convenience method
func (l Lexer) TokenizeString(s string) ([]Token, error) {
	r := bufio.NewReader(strings.NewReader(s))
	tokens := []Token{}
	err := l.Tokenize(r, func(t Token) error {
		tokens = append(tokens, t)
		return nil
	})
	return tokens, err
}

// AcceptsFilename returns true if this Lexer thinks it is suitable for the
// given filename. An error will be returned iff an invalid filename pattern
// is registered by the Lexer.
func (l Lexer) AcceptsFilename(name string) (bool, error) {
	for _, fn := range l.Filenames {
		if matched, err := path.Match(fn, name); err != nil {
			return false, fmt.Errorf("malformed filename pattern '%s' for "+
				"lexer '%s': %s", fn, l.Name, err)
		} else if matched {
			return true, nil
		}
	}
	return false, nil
}

// AcceptsMediaType returns true if this Lexer thinks it is suitable for the
// given meda (MIME) type. An error wil be returned iff the given mime type
// is invalid.
func (l Lexer) AcceptsMediaType(media string) (bool, error) {
	if mime, _, err := mime.ParseMediaType(media); err != nil {
		return false, err
	} else {
		for _, mt := range l.MimeTypes {
			if mime == mt {
				return true, nil
			}
		}
	}
	return false, nil
}

// ListMediaTypes lists the media types this Lexer supports,
// e.g. ["application/json"]
func (l Lexer) ListMediaTypes() []string {
	return l.MimeTypes
}

// ListFilenames lists the filename patterns this Lexer supports,
// e.g. ["*.json"]
func (l Lexer) ListFilenames() []string {
	return l.Filenames
}
