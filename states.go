package highlight

// States contains lexer states
type States interface {
	Get(name string) State
	Compile() (States, error)
}

// StatesSpec is a container for Lexer rule specifications, and can be
// compiled into a full state machine.
type StatesSpec map[string][]RuleSpec

func (m StatesSpec) Get(name string) State {
	return nil
}

// Compile compiles the specified states into a complete State machine,
// returning an error if any state fails to compile for any reason.
func (m StatesSpec) Compile() (States, error) {
	sm := &StateMap{}
	for name, specs := range m {
		rules := make(State, 0, len(specs))
		for _, spec := range specs {
			rules = append(rules, spec.Compile(sm))
		}
		(*sm)[name] = rules
	}
	return sm, nil
}

// MustCompile is a helper method that compiles the State specification,
// panicing on error.
func (m StatesSpec) MustCompile() States {
	states, err := m.Compile()
	if err != nil {
		panic(err)
	}
	return states
}

// StateMap is a map of states to their names.
type StateMap map[string]State

// Get returns the State with the given name.
func (m StateMap) Get(name string) State {
	return m[name]
}

// Compile does nothing.
func (m StateMap) Compile() (States, error) {
	return nil, nil
}

// State is a list of matching Rules.
type State []Rule

// Find examines the provided string, looking for a match within the current
// state. It returns the position `n` at which a rule match was found, and the
// rule itself.
//
// -1 will be returned if no rule could be matched, in which case
// the caller should disregard the string entirely (emit it as an error),
// and continue onto the next line of input.
//
// 0 will be returned if a rule matches at the start of the string.
//
// Otherwise, this function will return a number of characters to skip before
// reaching the first matched rule. The caller should emit those first `n`
// characters as an error, and emit the remaining characters according to the
// rule.
func (s State) Find(subject string) (int, Rule) {
	var earliestPos int = len(subject)
	var earliestRule Rule

	for _, rule := range s {
		pos, matchedRule := rule.Find(subject)

		if pos < 0 {
			// no match; try next rule
			continue
		} else if pos < earliestPos {
			earliestPos = pos
			earliestRule = matchedRule
		}
	}

	if earliestRule == nil {
		return -1, nil
	}

	return earliestPos, earliestRule
}

// Match tests the subject text against all rules within the State. If a match
// is found, it returns the number of characters consumed, a series of tokens
// consumed from the subject text, and the specific Rule that was succesfully
// matched against.
//
// If the start of the subject text can not be matched against any known rule,
// it will return a position of -1 and a nil Rule.
func (s State) Match(subject string) (int, Rule, []Token, error) {
	var earliestPos int = len(subject)
	var earliestRule Rule

	for _, rule := range s {
		pos, matchedRule := rule.Find(subject)

		if pos < 0 {
			// no match; try next rule
			continue
		} else if pos < earliestPos {
			earliestPos = pos
			earliestRule = matchedRule
		}
	}

	if earliestPos > 0 {
		// Return part of subject that doesn't match
		return earliestPos, nil, nil, nil
	} else if earliestRule == nil {
		return -1, nil, nil, nil
	}

	// Return matching part
	return earliestRule.Match(subject)
}
