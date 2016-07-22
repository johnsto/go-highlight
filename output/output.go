package output

import (
	"os"

	"github.com/johnsto/go-highlight"
)

type Outputter interface {
	highlight.Emitter
	SetFile(f *os.File) error
}
