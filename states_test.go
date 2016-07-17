package highlight_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "bitbucket.org/johnsto/go-highlight"
)

func TestStateMatch(t *testing.T) {
	s := State{
		NewRegexpRule("a+b+", String, nil, nil),
		NewRegexpRule("a+", String, nil, nil),
		NewRegexpRule("b+", String, nil, nil),
	}

	for _, item := range []struct {
		State   State
		Subject string
		Pos     int
		Rule    Rule
	}{
		{State{s[1]}, "", -1, nil},
		{State{s[1]}, "a", 0, s[1]},
		{State{s[1]}, "aaa", 0, s[1]},
		{State{s[1]}, "za", 1, s[1]},
		{State{s[1]}, "zzzabb", 3, s[1]},
		{s, "", -1, nil},
		{s, "a", 0, s[1]},
		{s, "aaa", 0, s[1]},
		{s, "za", 1, s[1]},
		{s, "zzzabb", 3, s[0]},
		{s, "b", 0, s[2]},
		{s, "zzzbbb", 3, s[2]},
		{s, "aaabbb", 0, s[0]},
		{s, "xxxaaabbb", 3, s[0]},
	} {
		n, rule := item.State.Find(item.Subject)
		assert.Equal(t, item.Pos, n, item.Subject)
		assert.Equal(t, item.Rule, rule, item.Subject)
	}
}
