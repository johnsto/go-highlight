package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"bitbucket.org/johnsto/go-highlight"
	_ "bitbucket.org/johnsto/go-highlight/lexers"
	"bitbucket.org/johnsto/go-highlight/output"
	"bitbucket.org/johnsto/go-highlight/output/term"
	"github.com/spf13/pflag"
)

func main() {
	contentType := pflag.StringP("type", "t", "", "content type to parse as")
	outputType := pflag.StringP("output", "o", "term", "output type [term]")

	pflag.Parse()
	filename := pflag.Arg(0)

	var emitter highlight.Emitter

	switch *outputType {
	case "term":
		emitter = term.NewOutput()
	case "debug":
		emitter = output.NewDebugOutputter()
	default:
		fmt.Printf("unknown output type '%s'. Valid values are:\n", *outputType)
		fmt.Println("  term - coloured terminal output")
		fmt.Println("  debug - debugging terminal output")
		os.Exit(1)
		return
	}

	var r io.Reader
	var tokenizer highlight.Tokenizer
	var err error

	// If user has specified a content type, resolve that first.
	if *contentType != "" {
		tokenizer, err = highlight.GetTokenizerForContentType(*contentType)
		if err != nil {
			fmt.Println("couldn't get tokenizer for content type:", err)
			os.Exit(1)
			return
		} else if tokenizer == nil {
			fmt.Printf("couldn't find tokenizer for content type '%s'\n", *contentType)
			os.Exit(1)
			return
		}
	}

	if filename != "" {
		// Read from specified file
		r, err = os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}

		// Get tokenizer for file extension
		if tokenizer == nil {
			tokenizer, err = highlight.GetTokenizerForFilename(path.Base(filename))
			if err != nil {
				fmt.Println("couldn't get tokenizer for file type:", err)
				os.Exit(1)
				return
			} else if tokenizer == nil {
				fmt.Println("couldn't find tokenizer for file type")
				os.Exit(1)
				return
			}
		}
	} else {
		// Read from stdin
		r = os.Stdin

		if tokenizer == nil {
			fmt.Println("no tokenizer specified - use `-t` to specify a " +
				"tokenizer when reading from standard input")
			os.Exit(1)
			return
		}
	}

	err = tokenizer.Tokenize(r, emitter.Emit)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
}
