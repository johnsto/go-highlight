package lexers

import . "github.com/johnsto/go-highlight"

var HTTP = Lexer{
	Name: "http",
	States: StatesSpec{
		"root": {
			// Request
			{Regexp: `^([a-zA-Z]+)( )([^ ]+)( )(HTTP)(/)([0-9\.]+)(\r\n)$`,
				SubTypes: []TokenType{Tag, Whitespace, String, Whitespace,
					Tag, Punctuation, Tag, Whitespace},
				State: "headers"},
			// Response
			{Regexp: `^(HTTP)(/)([0-9\.]+)( )([0-9]+)(.*)(\r\n)$`,
				SubTypes: []TokenType{Tag, Punctuation, Tag,
					Whitespace, Number, Whitespace, String, Whitespace},
				State: "headers"},
		},
		"headers": {
			{Regexp: `^(.*?)(:)(\s*)`,
				SubTypes: []TokenType{Attribute, Assignment, Whitespace},
				State:    "headerValue"},
			{Regexp: `^\r\n$`, State: "#pop #pop"},
		},
		"headerValue": {
			{Regexp: `\r\n$`, State: "#pop", Type: Whitespace},
			{Regexp: `[^;]+?`, Type: Text},
			{Regexp: `;`, Type: Punctuation},
			{Regexp: `\r\n$`, State: "#pop"},
		},
	},
	Filters: []Filter{},
}

func init() {
	Register(HTTP.Name, HTTP)
}
