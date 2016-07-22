package highlight_test

import (
	"testing"

	. "github.com/johnsto/go-highlight"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stack := Stack{}
	assert.Equal(t, "", stack.Peek(),
		"peeking on empty stack should return empty string")

	stack.Push("a")
	assert.Equal(t, "a", stack.Peek(),
		"peeking a stack should return the top element")

	stack.Push("b")
	assert.Equal(t, "b", stack.Peek(),
		"peeking a stack should return the top element")

	assert.Equal(t, "b", stack.Pop(),
		"popping a stack should return the top element")
	assert.Equal(t, "a", stack.Peek(),
		"peeking a popped stack should return the new top element")
	assert.Equal(t, "a", stack.Pop(),
		"popping a stack should return the top element")

	assert.Equal(t, "", stack.Peek(),
		"peeking on empty stack should return empty string")
	assert.Equal(t, "", stack.Pop(),
		"popping an empty stack should return empty string")

	stack.Push("a")
	stack.Push("b")
	stack.Push("c")
	stack.Empty()
	assert.Equal(t, "", stack.Peek())
}
