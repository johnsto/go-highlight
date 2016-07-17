package lexers

import . "bitbucket.org/johnsto/go-highlight"

var HTML = Lexer{
	Name:      "html",
	MimeTypes: []string{"text/html", "application/xhtml+xml"},
	Filenames: []string{"*.html", "*.htm", "*.xhtml"},
	States: StatesSpec{
		"root": {
			{Regexp: "[^<&]+", Type: Text},
			{Regexp: "&\\S+?;", Type: Entity},
			{Regexp: "<!--", Type: Comment, State: "comment"},
			{Regexp: "(<)(![^>]*)(>)",
				SubTypes: []TokenType{Punctuation, Entity, Punctuation}},
			{Regexp: "(</?)([\\w-]*:?[\\w-]+)(\\s*)(>)",
				SubTypes: []TokenType{Punctuation, Entity, Text, Punctuation}},
			{Regexp: "(<)([\\w-]*:?[\\w-]+)(\\s*)",
				SubTypes: []TokenType{Punctuation, Entity, Text},
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
			{Regexp: "\\s+", Type: Entity},
			{Regexp: "(/?)(\\s*)(>)",
				SubTypes: []TokenType{Punctuation, Entity, Punctuation},
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
