package output

import (
	"fmt"
	"io"
	"os"

	"github.com/johnsto/go-highlight"
)

type DebugOutputter struct {
	Writer io.Writer
}

func NewDebugOutputter() *DebugOutputter {
	return &DebugOutputter{
		Writer: os.Stdout,
	}
}

func (o *DebugOutputter) SetFile(f *os.File) error {
	o.Writer = f
	return nil
}

func (o *DebugOutputter) Emit(t highlight.Token) error {
	_, err := fmt.Fprintf(o.Writer,
		"%24s\t%12s\t%#v\n", t.State, t.Type, t.Value)
	return err
}
