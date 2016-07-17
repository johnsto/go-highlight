package lexers

import (
	. "bitbucket.org/johnsto/go-highlight"
	"strings"
)

var JSON = Lexer{
	Name:      "JSON",
	MimeTypes: []string{"application/json"},
	Filenames: []string{"*.json"},
	States: StatesSpec{
		"root": {
			{Include: "value"},
		},
		"whitespace": {
			{Regexp: "\\s+", Type: Whitespace},
		},
		"boolean": {
			{Regexp: "(true|false|null)", Type: Literal},
		},
		"number": {
			// -123.456e+78
			{Regexp: "-?[0-9]+\\.?[0-9]*[eE][\\+\\-]?[0-9]+", Type: Number},
			// -123.456
			{Regexp: "-?[0-9]+\\.[0-9]+", Type: Number},
			// -123
			{Regexp: "-?[0-9]+", Type: Number},
		},
		"string": {
			{Regexp: `(")(")`,
				SubTypes: []TokenType{Punctuation, Punctuation}},
			{Regexp: `(")((?:\\\"|[^\"])*?)(")`,
				SubTypes: []TokenType{Punctuation, String, Punctuation}},
		},
		"value": {
			{Include: "whitespace"},
			{Include: "boolean"},
			{Include: "number"},
			{Include: "string"},
			{Include: "array"},
			{Include: "object"},
		},
		"object": {
			{Regexp: "{", Type: Punctuation, State: "objectKey"},
		},
		"objectKey": {
			{Include: "whitespace"},
			{Regexp: `(")((?:\\\"|[^\"])*?)(")(\s*)(:)`,
				SubTypes: []TokenType{Punctuation, Attribute, Punctuation,
					Whitespace, Assignment},
				State: "objectValue"},
			{Regexp: "}", Type: Punctuation, State: "#pop"},
		},
		"objectValue": {
			{Include: "whitespace"},
			{Include: "value"},
			{Regexp: ",", Type: Punctuation, State: "#pop"},
			{Regexp: "}", Type: Punctuation, State: "#pop #pop"},
		},
		"array": {
			{Regexp: "\\[", Type: Punctuation, State: "arrayValue"},
		},
		"arrayValue": {
			{Include: "whitespace"},
			{Include: "value"},
			{Regexp: ",", Type: Punctuation},
			{Regexp: "\\]", Type: Punctuation, State: "#pop"},
		},
	},
	Filters: []Filter{
		RemoveEmptiesFilter,
		&JSONFormatter{Indent: "  "},
	},
}

type JSONFormatter struct {
	Indent string
}

func (f *JSONFormatter) Filter(lexer Lexer, emit func(Token) error) func(Token) error {

	//var laastState string
	indents := 0
	return func(token Token) error {
		indent := strings.Repeat(f.Indent, indents)

		var out []Token

		switch token.Type {
		case Whitespace:
			// we'll add our own whitespace, thanks!
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
