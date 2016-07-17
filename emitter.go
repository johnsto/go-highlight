package highlight

type Emitter interface {
	Emit(t Token) error
}
