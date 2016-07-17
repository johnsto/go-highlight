package highlight

type TokenType string

const (
	Error       TokenType = "error"
	Comment               = "comment"
	Text                  = "text"
	Number                = "number"
	String                = "string"
	Attribute             = "attribute"
	Assignment            = "assignment"
	Operator              = "operator"
	Punctuation           = "punctuation"
	Constant              = "constant"
	Entity                = "entity"
	Whitespace            = "whitespace"
)

var EndToken = Token{}
