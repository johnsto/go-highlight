package highlight

type TokenType string

const (
	// Error, emitted when unexpected token was encountered.
	Error TokenType = "error"
	// Comment e.g. `// this should never happen`
	Comment = "comment"
	// Number - e.g. `2716057` in `"serial": 2716057` or `serial = 2716057;`
	Number = "number"
	// String - e.g. `Fry` in `"name": "Fry"` or `var name = "Fry";`
	String = "string"
	// Text - e.g. `Fry` in `<p>Fry</p>`
	Text = "text"
	// Attribute - e.g. `name` in `"name": "Fry"`, or `font-size` in
	// `font-size: 1.2rem;`
	Attribute = "attribute"
	// Assignment - e.g. `=` in `int x = y;` or `:` in `font-size: 1.2rem;`
	Assignment = "assignment"
	// Operator - e.g. `+`/`-` in `int x = a + b - c;`
	Operator = "operator"
	// Punctuation - e.g. semi/colons in `int x, j;`
	Punctuation = "punctuation"
	// Literal - e.g. `true`/`false`/`null`.
	Literal = "literal"
	// Tag - e.g. `html`/`div`/`b`
	Tag = "tag"
	// Whitespace - e.g. \n, \t
	Whitespace = "whitespace"
)

var EndToken = Token{}
