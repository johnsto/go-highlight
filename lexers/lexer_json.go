package lexers

import (
	. "bitbucket.org/johnsto/go-highlight"
	"strings"
)

var JSON = Lexer{
	Name:      "json",
	MimeTypes: []string{"application/json"},
	Filenames: []string{"*.json"},
	States: StatesSpec{
		"root": {
			{Include: "value"},
		},
		"whitespace": {
			{Regexp: "\\s+", Type: Whitespace},
		},
		// literal matches a literal JSON value
		"literal": {
			{Regexp: "(true|false|null)", Type: Literal},
		},
		// number matches a JSON number
		"number": {
			// -123.456e+78
			{Regexp: "-?[0-9]+\\.?[0-9]*[eE][\\+\\-]?[0-9]+", Type: Number},
			// -123.456
			{Regexp: "-?[0-9]+\\.[0-9]+", Type: Number},
			// -123
			{Regexp: "-?[0-9]+", Type: Number},
		},
		// string matches a JSON string
		"string": {
			{Regexp: `(")(")`,
				SubTypes: []TokenType{Punctuation, Punctuation}},
			{Regexp: `(")((?:\\\"|[^\"])*?)(")`,
				SubTypes: []TokenType{Punctuation, String, Punctuation}},
		},
		// value matches any valid JSON value
		"value": {
			{Include: "whitespace"},
			{Include: "literal"},
			{Include: "number"},
			{Include: "string"},
			{Include: "array"},
			{Include: "object"},
		},
		// object matches the start of an object
		"object": {
			{Regexp: "{", Type: Punctuation, State: "objectKey"},
		},
		// objectKey matches a key within an object, or pops if the end of
		// the object has been reached
		"objectKey": {
			{Include: "whitespace"},
			{Regexp: `(")((?:\\\"|[^\"])*?)(")(\s*)(:)`,
				SubTypes: []TokenType{Punctuation, Attribute, Punctuation,
					Whitespace, Assignment},
				State: "objectValue"},
			{Regexp: "}", Type: Punctuation, State: "#pop"},
		},
		// objectValue matches a key value within an object, popping after
		// each element or when the object ends
		"objectValue": {
			{Include: "whitespace"},
			{Include: "value"},
			{Regexp: ",", Type: Punctuation, State: "#pop"},
			{Regexp: "}", Type: Punctuation, State: "#pop #pop"},
		},
		// array matches the start of an array
		"array": {
			{Regexp: "\\[", Type: Punctuation, State: "arrayValue"},
		},
		// arrayValue matches elements within an array and pops when the
		// array ends
		"arrayValue": {
			{Include: "whitespace"},
			{Include: "value"},
			{Regexp: ",", Type: Punctuation},
			{Regexp: "\\]", Type: Punctuation, State: "#pop"},
		},
	},
	Filters: []Filter{
		RemoveEmptiesFilter,
	},
	Formatter: &JSONFormatter{Indent: "  "},
}

// JSONFormatter consumes a series of JSON tokens and emits additional tokens
// to produce indented, formatted output.
type JSONFormatter struct {
	Indent string
}

func (f *JSONFormatter) Filter(emit func(Token) error) func(
	Token) error {

	// indents records the current indentation level
	indents := 0

	return func(token Token) error {
		indent := strings.Repeat(f.Indent, indents)

		// temporary storage for the tokens to emit
		var out []Token

		switch token.Type {
		case Whitespace:
			// nah, we'll add our own whitespace, thanks!
			return nil
		case Assignment:
			switch token.Value {
			case ":":
				out = []Token{token, Token{Type: Whitespace, Value: " "}}
			default:
				out = []Token{token}
			}
		case Punctuation:
			switch token.Value {
			case ",":
				out = []Token{token,
					Token{Type: Whitespace, Value: "\n"},
					Token{Type: Whitespace, Value: indent}}
			case "{":
				fallthrough
			case "[":
				out = append(out, token)
				out = append(out, Token{Type: Whitespace, Value: "\n"})
				indents++
				indent = strings.Repeat(f.Indent, indents)
				out = append(out, Token{Type: Whitespace, Value: indent})
			case "}":
				fallthrough
			case "]":
				out = append(out, Token{Type: Whitespace, Value: "\n"})
				indents--
				indent = strings.Repeat(f.Indent, indents)
				out = append(out, Token{Type: Whitespace, Value: indent})
				out = append(out, token)
			case "\"":
				out = []Token{token}
			default:
				out = []Token{token}
			}
		case "":
			// EOF
			break
		default:
			out = []Token{token}
		}

		// Attempt to emit each token, failing on first error
		for _, t := range out {
			if err := emit(t); err != nil {
				return err
			}
		}

		return nil
	}
	return nil
}

func init() {
	Register(JSON.Name, JSON)
}
