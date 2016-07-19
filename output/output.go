package output

import (
	"os"

	"bitbucket.org/johnsto/go-highlight"
)

type Outputter interface {
	highlight.Emitter
	SetFile(f *os.File) error
}
