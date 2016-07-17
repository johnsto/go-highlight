package lexers_test

import (
	"fmt"
	"io"
	"testing"

	. "bitbucket.org/johnsto/go-highlight"
	"bitbucket.org/johnsto/go-highlight/lexers"
	"github.com/stretchr/testify/assert"
)

func TestLexerJSONSimple(t *testing.T) {
	type simpleToken struct {
		Value string
		Type  TokenType
	}

	states, err := lexers.JSON.States.Compile()
	assert.Nil(t, err, "JSON lexer should compile")

	for _, item := range []struct {
		State   string
		Length  int
		Subject string
		Tokens  []Token
	}{
		{"boolean", 4, "true", []Token{{Value: "true", Type: Constant}}},
		{"boolean", 5, "false", []Token{{Value: "false", Type: Constant}}},
		{"boolean", 4, "null", []Token{{Value: "null", Type: Constant}}},
		{"number", 1, "0", []Token{{Value: "0", Type: Number}}},
		{"number", 2, "-0", []Token{{Value: "-0", Type: Number}}},
		{"number", 3, "0.0", []Token{{Value: "0.0", Type: Number}}},
		{"number", 4, "-0.0", []Token{{Value: "-0.0", Type: Number}}},
		{"number", 5, "1.2e3", []Token{{Value: "1.2e3", Type: Number}}},
		{"number", 6, "-1.2e3", []Token{{Value: "-1.2e3", Type: Number}}},
		{"number", 7, "-1.2e-4", []Token{{Value: "-1.2e-4", Type: Number}}},
		{"string", 2, `""`, []Token{
			{Value: `"`, Type: Punctuation},
			{Value: `"`, Type: Punctuation},
		}},
		{"string", 4, `"  "`, []Token{
			{Value: `"`, Type: Punctuation},
			{Value: `  `, Type: String},
			{Value: `"`, Type: Punctuation},
		}},
		{"string", 5, `"xyz"`, []Token{
			{Value: `"`, Type: Punctuation},
			{Value: `xyz`, Type: String},
			{Value: `"`, Type: Punctuation},
		}},
		{"string", 10, `"\"cats\""`, []Token{
			{Value: `"`, Type: Punctuation},
			{Value: `\"cats\"`, Type: String},
			{Value: `"`, Type: Punctuation},
		}},
		{"string", 26, `"escape\tall\nthe\rthings"`, []Token{
			{Value: `"`, Type: Punctuation},
			{Value: `escape\tall\nthe\rthings`, Type: String},
			{Value: `"`, Type: Punctuation},
		}},
	} {
		n, _, tokens, err := states.Get(item.State).Match(item.Subject)
		assert.Nil(t, err, item.Subject)
		assert.Equal(t, item.Length, n, item.Subject)
		assert.Equal(t, item.Tokens, tokens, item.Subject)
	}
}

func TestLexerJSONComplex(t *testing.T) {
	type simpleToken struct {
		Value string
		Type  TokenType
	}

	_, err := lexers.JSON.States.Compile()
	assert.Nil(t, err, "JSON lexer should compile")
	for _, item := range []struct {
		State   string
		Subject string
		Tokens  []simpleToken
	}{
		{"array", "[]", []simpleToken{
			{"[", Punctuation},
			{"]", Punctuation},
		}},
		{"array", "[ ]", []simpleToken{
			{"[", Punctuation},
			{" ", Whitespace},
			{"]", Punctuation},
		}},
		{"array", "[\n]", []simpleToken{
			{"[", Punctuation},
			{"\n", Whitespace},
			{"]", Punctuation},
		}},
		{"array", "[null]", []simpleToken{
			{"[", Punctuation},
			{"null", Constant},
			{"]", Punctuation},
		}},
		{"array", "[123]", []simpleToken{
			{"[", Punctuation},
			{"123", Number},
			{"]", Punctuation},
		}},
		{"array", "[1,2]", []simpleToken{
			{"[", Punctuation},
			{"1", Number},
			{",", Punctuation},
			{"2", Number},
			{"]", Punctuation},
		}},
		{"array", "[1,\n2]", []simpleToken{
			{"[", Punctuation},
			{"1", Number},
			{",", Punctuation},
			{"\n", Whitespace},
			{"2", Number},
			{"]", Punctuation},
		}},
		{"object", "{}", []simpleToken{
			{"{", Punctuation},
			{"}", Punctuation},
		}},
		{"object", `{"key":"value"}`, []simpleToken{
			{"{", Punctuation},
			{`"`, Punctuation},
			{"key", Attribute},
			{`"`, Punctuation},
			{"", Whitespace},
			{":", Assignment},
			{`"`, Punctuation},
			{"value", String},
			{`"`, Punctuation},
			{"}", Punctuation},
		}},
		{"object", "{\"key\":\n\"value\"}", []simpleToken{
			{"{", Punctuation},
			{`"`, Punctuation},
			{"key", Attribute},
			{`"`, Punctuation},
			{"", Whitespace},
			{":", Assignment},
			{"\n", Whitespace},
			{`"`, Punctuation},
			{"value", String},
			{`"`, Punctuation},
			{"}", Punctuation},
		}},
		{"object", `{"ke\ty":"v\nalu\re\""}`, []simpleToken{
			{`{`, Punctuation},
			{`"`, Punctuation},
			{`ke\ty`, Attribute},
			{`"`, Punctuation},
			{``, Whitespace},
			{`:`, Assignment},
			{`"`, Punctuation},
			{`v\nalu\re\"`, String},
			{`"`, Punctuation},
			{`}`, Punctuation},
		}},
		{"object", `{ "key" : "value" }`, []simpleToken{
			{"{", Punctuation},
			{" ", Whitespace},
			{`"`, Punctuation},
			{"key", Attribute},
			{`"`, Punctuation},
			{" ", Whitespace},
			{":", Assignment},
			{" ", Whitespace},
			{`"`, Punctuation},
			{"value", String},
			{`"`, Punctuation},
			{" ", Whitespace},
			{"}", Punctuation},
		}},
		{"object", `{"aa":"bb","cc":"dd"}`, []simpleToken{
			{`{`, Punctuation},
			{`"`, Punctuation},
			{`aa`, Attribute},
			{`"`, Punctuation},
			{``, Whitespace},
			{`:`, Assignment},
			{`"`, Punctuation},
			{`bb`, String},
			{`"`, Punctuation},
			{`,`, Punctuation},
			{`"`, Punctuation},
			{`cc`, Attribute},
			{`"`, Punctuation},
			{``, Whitespace},
			{`:`, Assignment},
			{`"`, Punctuation},
			{`dd`, String},
			{`"`, Punctuation},
			{`}`, Punctuation},
		}},
		{"object", "{\"key\":[1,\"value\",{\"a\":\"b\"}]}", []simpleToken{
			{"{", Punctuation},
			{`"`, Punctuation},
			{`key`, Attribute},
			{`"`, Punctuation},
			{``, Whitespace},
			{`:`, Assignment},
			{"[", Punctuation},
			{"1", Number},
			{",", Punctuation},
			{`"`, Punctuation},
			{`value`, String},
			{`"`, Punctuation},
			{",", Punctuation},
			{"{", Punctuation},
			{`"`, Punctuation},
			{"a", Attribute},
			{`"`, Punctuation},
			{"", Whitespace},
			{":", Assignment},
			{`"`, Punctuation},
			{"b", String},
			{`"`, Punctuation},
			{"}", Punctuation},
			{"]", Punctuation},
			{"}", Punctuation},
		}},
	} {
		tokens, err := lexers.JSON.TokenizeString(item.Subject)
		tokens = tokens[0 : len(tokens)-1] // remove EndToken
		name := fmt.Sprintf("`%s` %v %v", item.Subject, tokens, err)
		assert.Equal(t, io.EOF, err,
			fmt.Sprintf("tokeniser should return EOF"))
		if !assert.Equal(t, len(item.Tokens), len(tokens),
			fmt.Sprintf("number of tokens in %#v should match", item.Subject)) {
			fmt.Println(tokens)
			fmt.Println(item.Tokens)
			fmt.Println()
		}
		for i, token := range tokens {
			if i >= len(item.Tokens) {
				break
			}
			actualToken := Token{
				Value: token.Value,
				Type:  token.Type,
			}
			expectedToken := Token{
				Value: item.Tokens[i].Value,
				Type:  item.Tokens[i].Type,
			}
			if i < len(item.Tokens) {
				assert.Equal(t, expectedToken, actualToken,
					fmt.Sprintf("(%d) %s", i, name))
			}
		}
	}
}
