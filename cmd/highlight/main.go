package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/johnsto/go-highlight"
	_ "github.com/johnsto/go-highlight/lexers"
	"github.com/johnsto/go-highlight/output"
	"github.com/johnsto/go-highlight/output/term"
	"github.com/spf13/pflag"
)

func listSupportedTypes(w io.Writer) {
	tokenizers := highlight.GetTokenizers()

	fmt.Fprintln(w, "Registered media types:")
	for name, tokenizer := range tokenizers {
		types := tokenizer.ListMediaTypes()
		fmt.Fprintf(w, "  %s: %q\n", name, types)
	}

	fmt.Fprintln(w, "\nRegistered file patterns:")
	for name, tokenizer := range tokenizers {
		patterns := tokenizer.ListFilenames()
		fmt.Fprintf(w, "  %s: %q\n", name, patterns)
	}
}

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options...] [file]\n"+
				"Syntax-highlights file to standard output\n\n",
			os.Args[0])
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n"+
			"Full documentation at: <https://github.com/johnsto/go-highlight>\n")
	}

	contentType := pflag.StringP("type", "t", "",
		"content type to parse as (e.g. 'application/json')")
	outputType := pflag.StringP("output", "o", "ansi", "output type [ansi] "+
		"<ansi|text|debug>")
	outputFile := pflag.StringP("output-file", "O", "", "output to file")
	listSupported := pflag.BoolP("list", "l", false, "list supported types")

	pflag.Parse()

	if *listSupported {
		listSupportedTypes(os.Stdout)
		os.Exit(0)
	}

	filename := pflag.Arg(0)

	var outputter output.Outputter

	// Determine output style
	switch *outputType {
	case "ansi":
		outputter = term.NewOutput()
	case "text":
		outputter = output.NewTextOutputter()
	case "debug":
		outputter = output.NewDebugOutputter()
	default:
		fmt.Fprintf(os.Stderr,
			"unknown output type '%s'. Valid values are:\n"+
				"  ansi - coloured ANSI output"+
				"  text - standard text output"+
				"  debug - debugging output",
			*outputType)
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
			fmt.Fprintf(os.Stderr,
				"couldn't get tokenizer for content type: %s\n", err)
			os.Exit(1)
			return
		} else if tokenizer == nil {
			fmt.Fprintf(os.Stderr,
				"couldn't find tokenizer for content type '%s'\n",
				*contentType)
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
			tokenizer, err = highlight.GetTokenizerForFilename(
				path.Base(filename))
			if err != nil {
				fmt.Fprintln(os.Stderr,
					"couldn't get tokenizer for file type:", err)
				listSupportedTypes(os.Stderr)
				os.Exit(1)
				return
			} else if tokenizer == nil {
				fmt.Fprintln(os.Stderr,
					"couldn't find tokenizer for file type")
				listSupportedTypes(os.Stderr)
				os.Exit(1)
				return
			}
		}
	} else {
		// Read from stdin
		r = os.Stdin

		if tokenizer == nil {
			fmt.Fprintln(os.Stderr,
				"no tokenizer specified - use `-t` to specify a "+
					"tokenizer when reading from standard input")
			os.Exit(1)
			return
		}
	}

	// Write output to file if specified
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalln(err)
		}
		outputter.SetFile(f)
	}

	err = tokenizer.Format(r, outputter.Emit)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
}
