package output

import (
	"fmt"
	"io"
	"os"

	"github.com/johnsto/go-highlight"
)

type TextOutputter struct {
	Writer io.Writer
}

func NewTextOutputter() *TextOutputter {
	return &TextOutputter{
		Writer: os.Stdout,
	}
}

func (o *TextOutputter) SetFile(f *os.File) error {
	o.Writer = f
	return nil
}

func (o *TextOutputter) Emit(t highlight.Token) error {
	_, err := fmt.Fprintf(o.Writer, t.Value)
	return err
}
