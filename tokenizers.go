package highlight

var tokenizers map[string]Tokenizer

// Register registers the given Tokenizer under the specified name. Any
// existing Tokenizer under that name will be replaced.
func Register(name string, t Tokenizer) {
	if tokenizers == nil {
		tokenizers = map[string]Tokenizer{}
	}
	tokenizers[name] = t
}

// GetTokenizer returns the Tokenizer of the given name.
func GetTokenizer(name string) Tokenizer {
	return tokenizers[name]
}

// GetTokenizers returns the map of known Tokenizers.
func GetTokenizers() map[string]Tokenizer {
	return tokenizers
}

// GetTokenizerForContentType returns a Tokenizer for the given content type
// (e.g. "text/html" or "application/json"), or nil if one is not found.
func GetTokenizerForContentType(contentType string) (Tokenizer, error) {
	for _, tokenizer := range tokenizers {
		if matched, err := tokenizer.AcceptsMediaType(contentType); err != nil {
			return nil, err
		} else if matched {
			return tokenizer, nil
		}
	}
	return nil, nil
}

// GetTokenizerForFilename returns a Tokenizer for the given filename
// (e.g. "index.html" or "jasons.json"), or nil if one is not found.
func GetTokenizerForFilename(name string) (Tokenizer, error) {
	for _, tokenizer := range tokenizers {
		if matched, err := tokenizer.AcceptsFilename(name); err != nil {
			return nil, err
		} else if matched {
			return tokenizer, nil
		}
	}
	return nil, nil
}
