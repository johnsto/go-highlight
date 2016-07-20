package highlight

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Rule interface {
		Find(subject string) (int, Rule)
		Match(subject string) (int, Rule, []Token, error)
		Stack() []string
	}

	// Rule describes the conditions required to match some subject text.
	RuleSpec struct {
		// Regexp is the regular expression this rule should match against.
		Regexp string
		// Type is the token type for strings that match this rule.
		Type TokenType
		// SubTypes contains an ordered array of token types matching the order
		// of groups in the Regexp expression.
		SubTypes []TokenType
		// State indicates the next state to migrate to if this rule is
		// triggered.
		State string
		// Include specifies a state to run
		Include string
	}

	// IncludeRule allows the states of another Rule to be referenced.
	IncludeRule struct {
		StateMap  *StateMap
		StateName string
	}

	// RegexpRule matches a state if the subject matches a regular expression.
	RegexpRule struct {
		Regexp     *regexp.Regexp
		Type       TokenType
		SubTypes   []TokenType
		NextStates []string
	}
)

// Compile converts the RuleSpec shorthand into a fully-fledged Rule.
func (rs RuleSpec) Compile(sm *StateMap) Rule {
	if rs.Include != "" {
		return IncludeRule{
			StateMap:  sm,
			StateName: rs.Include,
		}
	}
	return NewRegexpRule(rs.Regexp, rs.Type, rs.SubTypes,
		strings.Split(rs.State, " "))
}

// NewRegexpRule creates a new regular expression Rule.
func NewRegexpRule(re string, t TokenType, subTypes []TokenType,
	next []string) RegexpRule {
	return RegexpRule{
		Regexp:     regexp.MustCompile(re),
		Type:       t,
		SubTypes:   subTypes,
		NextStates: next,
	}
}

// Find returns the first position in subject where this Rule will
// match, or -1 if no match was found.
func (r RegexpRule) Find(subject string) (int, Rule) {
	if indices := r.Regexp.FindStringIndex(subject); indices == nil {
		return -1, nil
	} else {
		return indices[0], r
	}
}

// Match attempts to match against the beginning of the given search string.
// Returns the number of characters matched, and an array of tokens.
//
// If the regular expression contains groups, they will be matched with the
// corresponding token type in `Rule.Types`. Any text inbetween groups will
// be returned using the token type defined by `Rule.Type`.
func (r RegexpRule) Match(subject string) (int, Rule, []Token, error) {
	// Find match group and sub groups, returns an array of start/end offsets
	// e.g. f(r/a(b+)c/g, "abbbc") = [0, 5, 1, 4]
	indices := r.Regexp.FindStringSubmatchIndex(subject)

	if indices == nil || indices[0] != 0 || indices[1] == 0 {
		// Didn't match start of subject
		return -1, nil, nil, nil
	}

	// Get position after final matched character
	n := indices[1]

	if r.SubTypes == nil {
		// No groups in expression; return single token and type
		return n, r, []Token{{
			Value: subject[:n],
			Type:  r.Type,
		}}, nil
	}

	// Multiple groups; construct array of group values and tokens
	tokens := []Token{}
	var start, end int = 0, 0
	for i := 2; i < len(indices); i += 2 {
		prevEnd := end
		start, end = indices[i], indices[i+1]

		if start < 0 || end < 0 {
			// Ignore empty submatch
			end = prevEnd
			continue
		}

		// Extract text between submatches
		sep := subject[prevEnd:start]
		if sep != "" {
			// Append to output
			tokens = append(tokens, Token{
				Value: sep,
				Type:  r.Type,
			})
		}

		// Determine submatch token
		j := (i - 2) / 2
		if j >= len(r.SubTypes) {
			return n, r, tokens, fmt.Errorf("not enough subtypes for group")
		}
		tokenType := r.SubTypes[j]

		// Extract submatch text
		tokens = append(tokens, Token{
			Value: subject[start:end],
			Type:  tokenType,
		})
	}

	return n, r, tokens, nil
}

func (r RegexpRule) Stack() []string {
	return r.NextStates
}

func (r IncludeRule) Find(subject string) (int, Rule) {
	state := r.StateMap.Get(r.StateName)
	return state.Find(subject)
}

func (r IncludeRule) Match(subject string) (int, Rule, []Token, error) {
	state := r.StateMap.Get(r.StateName)
	n, rl, ts, err := state.Match(subject)
	// set `State` property of each Token so they show the actual State.
	for _, t := range ts {
		t.State = r.StateName
	}
	return n, rl, ts, err
}

func (r IncludeRule) Stack() []string {
	return nil
}
