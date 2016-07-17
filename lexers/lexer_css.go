package lexers

import . "bitbucket.org/johnsto/go-highlight"

var CSS = Lexer{
	Name:      "CSS",
	MimeTypes: []string{"text/css"},
	Filenames: []string{"*.css"},
	States: StatesSpec{
		"root": {
			{Include: "whitespace"},
			{Include: "singleLineComment"},
			{Include: "multiLineComment"},
			{Include: "selector"},
			{Include: "declarationBlock"},
		},
		"selector": {
			{Regexp: `(\[)([^\]]+)(\])`,
				SubTypes: []TokenType{Punctuation, Attribute, Punctuation}},
			{Regexp: `(\.)([-a-zA-Z0-9]+)`,
				SubTypes: []TokenType{Punctuation, Attribute}},
			{Regexp: `@[-a-zA-Z0-9]+`, Type: Literal, State: "media"},
			{Regexp: `>`, Type: Punctuation},
			{Regexp: `\+`, Type: Punctuation},
			{Regexp: `:`, Type: Punctuation},
			{Regexp: `,`, Type: Punctuation},
			{Regexp: `[-a-zA-Z0-9]+`, Type: Attribute},
			{Regexp: `\*`, Type: Attribute},
		},
		"media": {
			{Regexp: ` and `, Type: Operator},
			{Regexp: `,`, Type: Punctuation},
			{Regexp: `[-a-zA-Z0-9]+`, Type: Attribute},
			{Regexp: `(\()` + `(\s*)` +
				`([-a-zA-Z0-9]+)` + `(:)` + `([^\)]+)` +
				`(\s*)` + `(\))`,
				SubTypes: []TokenType{Punctuation, Whitespace, Attribute, Assignment, Text, Whitespace, Punctuation}},
			{Include: "whitespace"},
			{Include: "singleLineComment"},
			{Include: "multiLineComment"},
			{Regexp: `{`, Type: Punctuation, State: "mediaBlock"},
		},
		"mediaBlock": {
			{Include: "whitespace"},
			{Include: "singleLineComment"},
			{Include: "multiLineComment"},
			{Include: "selector"},
			{Include: "declarationBlock"},
			{Regexp: `}`, Type: Punctuation, State: "#pop #pop"},
		},
		"ruleValue": {
			{Regexp: `;`, Type: "Punctuation", State: "#pop"},
			{Regexp: `.*`, Type: "Text"},
		},
		"declarationBlock": {
			{Regexp: `{`, Type: Punctuation, State: "declaration"},
		},
		"declaration": {
			{Include: "whitespace"},
			{Include: "singleLineComment"},
			{Include: "multiLineComment"},
			{Regexp: `([a-zA-Z0-9_-]+)(\w*)(:)`,
				SubTypes: []TokenType{Tag, Whitespace, Assignment},
				State:    "declarationValue"},
			{Regexp: `}`, Type: Punctuation, State: "#pop"},
			{Include: "selector"},
			{Include: "declarationBlock"},
		},
		"declarationValue": {
			{Regexp: `(")([^"]*)(")`,
				SubTypes: []TokenType{Punctuation, Text, Punctuation}},
			{Regexp: `(')([^']*)(')`,
				SubTypes: []TokenType{Punctuation, Text, Punctuation}},
			{Regexp: `[^;]+`, Type: Text},
			{Regexp: `,`, Type: Punctuation},
			{Regexp: `;`, Type: Punctuation, State: "#pop"},
		},
		"whitespace": {
			{Regexp: `[ \r\n\f\t]+`, Type: Whitespace},
		},
		"singleLineComment": {
			{Regexp: `\/\/.*`, Type: Comment},
		},
		"multiLineComment": {
			{Regexp: `\/\*`, Type: Comment, State: "multiLineCommentContents"},
		},
		"multiLineCommentContents": {
			{Regexp: `\*\/`, Type: Comment, State: "#pop"},
			{Regexp: `(.+?)(\*\/)`, Type: Comment, State: "#pop"},
			{Regexp: `.+`, Type: Comment},
		},
	},
}

func init() {
	Register(CSS.Name, CSS)
}
