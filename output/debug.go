package output

import (
	"fmt"

	"bitbucket.org/johnsto/go-highlight"
)

type DebugOutputter struct {
}

func NewDebugOutputter() *DebugOutputter {
	return &DebugOutputter{}
}

func (o *DebugOutputter) Emit(t highlight.Token) error {
	_, err := fmt.Printf("%24s\t%12s\t%#v\n", t.State, t.Type, t.Value)
	return err
}
