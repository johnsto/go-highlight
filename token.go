package highlight

import "fmt"

// Token represents one item of parsed output, containing the parsed value
// and its detected type.
type Token struct {
	Value string
	Type  TokenType
	State string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s:%#v {%s})", t.Type, t.Value, t.State)
}
