package highlight

// Emitter is any type that supports emitting tokens to some output
type Emitter interface {
	// Emit emits the given token to some output
	Emit(t Token) error
}
