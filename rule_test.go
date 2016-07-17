package highlight_test

import (
	"fmt"
	"testing"

	. "bitbucket.org/johnsto/go-highlight"
	"github.com/stretchr/testify/assert"
)

func TestRegexpRuleFind(t *testing.T) {
	for _, item := range []struct {
		Regexp   string
		Subject  string
		Position int
	}{
		{"a+b+c+", "aabbcc", 0},
		{"a+b+c+", "zzzaabbcc", 3},
		{"a+b+c+", "zzzaab", -1},
	} {
		rule := NewRegexpRule(item.Regexp, "", nil, nil)
		pos, r := rule.Find(item.Subject)
		if item.Position < 0 {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, rule, r)
		}
		assert.Equal(t, item.Position, pos)
	}
	//rule := RuleSpec{Regexp: "a+b+c+"}.Compile(nil)
	//pos, r := rule.Find("aabbcc")
	//assert.Equal(t, 0, pos)
	//assert.Equal(t, 3, rule.Find("zzzaabbcc"))
	//assert.Equal(t, -1, rule.Find("zzz"))
}

func TestRegexpRuleMatch(t *testing.T) {
	for _, item := range []struct {
		Regexp  string
		Type    TokenType
		Types   []TokenType
		Subject string
		Length  int
		Tokens  []Token
	}{
		// Non-matching
		{"ab+c", Text, nil, "", -1, nil},
		// Simple matching
		{"ab+c", Text, nil, "abc", 3, []Token{{"abc", Text, ""}}},
		{"ab+c", Text, nil, "abbbc", 5, []Token{{"abbbc", Text, ""}}},
		// Non-matching subgroup
		{"(b+)(c+)", Error, []TokenType{Text}, "bbb", -1, nil},
		// Simple matching subgroup
		{"(b+)(c+)", Error, []TokenType{Text, Text}, "bbcc", 4,
			[]Token{{"bb", Text, ""}, {"cc", Text, ""}}},
		{"(b+)(c+)", Error, nil, "bbcc", 4, []Token{{"bbcc", Error, ""}}},
		// Subgroup with outliers
		{"a(b+)cc(d+)", Error, []TokenType{Text, Text}, "abbccddd", 8,
			[]Token{{"a", Error, ""}, {"bb", Text, ""},
				{"cc", Error, ""}, {"ddd", Text, ""}}},
	} {
		rule := NewRegexpRule(item.Regexp, item.Type, item.Types, nil)
		n, _, tokens, err := rule.Match(item.Subject)
		description := fmt.Sprintf("%s - %s", item.Regexp, item.Subject)
		assert.Nil(t, err, description)
		assert.Equal(t, item.Length, n, description)
		assert.Equal(t, len(item.Tokens), len(tokens), description)
		assert.Equal(t, item.Tokens, tokens, description)
	}
}

func TestIncludeRule(t *testing.T) {
	sm := StateMap{}
	sm["root"] = State{
		IncludeRule{StateMap: &sm, StateName: "include1"},
		IncludeRule{StateMap: &sm, StateName: "include2"},
		NewRegexpRule("c+", "", nil, nil),
		NewRegexpRule("d+", "", nil, nil), // shouldn't match from root
	}
	sm["include1"] = State{
		NewRegexpRule("a+", "", nil, nil),
		NewRegexpRule("d+", "", nil, nil), // *should* match from root
	}
	sm["include2"] = State{
		NewRegexpRule("b+", "", nil, nil),
		NewRegexpRule("a+", "", nil, nil), // shouldn't match from root
		NewRegexpRule("d+", "", nil, nil), // shouldn't match from root
	}

	for _, item := range []struct {
		State    string
		Subject  string
		Position int
		Rule     Rule
	}{
		{"root", "aaa", 0, sm["include1"][0]},
		{"root", "bbb", 0, sm["include2"][0]},
		{"root", "ccc", 0, sm["root"][2]},
		{"root", "ddd", 0, sm["include1"][1]},
		{"root", "eee", -1, nil},
		{"root", "eeeaaa", 3, sm["include2"][1]},
	} {
		description := fmt.Sprintf("%s - %s", item.State, item.Subject)
		pos, rule := sm.Get(item.State).Find(item.Subject)
		assert.Equal(t, item.Position, pos, description)
		assert.Equal(t, item.Rule, rule, description)
	}

}
