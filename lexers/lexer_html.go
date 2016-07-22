package lexers

import . "github.com/johnsto/go-highlight"

var HTML = Lexer{
	Name:      "html",
	MimeTypes: []string{"text/html", "application/xhtml+xml"},
	Filenames: []string{"*.html", "*.htm", "*.xhtml"},
	States: StatesSpec{
		"root": {
			{Regexp: "[^<&]+", Type: Text},
			{Regexp: "&\\S+?;", Type: Tag},
			{Regexp: "<!--", Type: Comment, State: "comment"},
			{Regexp: "(<)(![^>]*)(>)",
				SubTypes: []TokenType{Punctuation, Tag, Punctuation}},
			{Regexp: "(</?)([\\w-]*:?[\\w-]+)(\\s*)(>)",
				SubTypes: []TokenType{Punctuation, Tag, Text, Punctuation}},
			{Regexp: "(<)([\\w-]*:?[\\w-]+)(\\s*)",
				SubTypes: []TokenType{Punctuation, Tag, Text},
				State:    "tag"},
		},
		"comment": {
			{Regexp: "-->", Type: Comment, State: "#pop"},
			{Regexp: "[^-]+", Type: Comment},
		},
		"tag": {
			{Regexp: "([\\w-]+)(=)(\\s*)",
				SubTypes: []TokenType{Attribute, Operator, Text},
				State:    "tagAttr"},
			{Regexp: "[\\w-]+\\s*", Type: Attribute},
			{Regexp: "\\s+", Type: Tag},
			{Regexp: "(/?)(\\s*)(>)",
				SubTypes: []TokenType{Punctuation, Tag, Punctuation},
				State:    "#pop",
			},
		},
		"tagAttr": {
			{Regexp: "\"[^\"]*\"", Type: String, State: "#pop"},
			{Regexp: "'[^']*'", Type: String, State: "#pop"},
			{Regexp: "\\w+", Type: String, State: "#pop"},
		},
	},
}

func init() {
	Register(HTML.Name, HTML)
}
