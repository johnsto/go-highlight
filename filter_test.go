package highlight_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/johnsto/go-highlight"
)

func TestFilters(t *testing.T) {
	for _, item := range []struct {
		Name    string
		Filters Filters
		Input   []Token
		Output  []Token
		Error   error
	}{{
		Name:    "no filters",
		Filters: Filters{},
		Input: []Token{
			{Value: "a", Type: Text},
			{Value: "b", Type: Text},
			{Value: "c", Type: Text},
		},
		Output: []Token{
			{Value: "a", Type: Text},
			{Value: "b", Type: Text},
			{Value: "c", Type: Text},
		},
	}, {
		Name: "error",
		Filters: Filters{FilterFunc(func(out func(Token) error) func(Token) error {
			return func(t Token) error {
				return io.EOF
			}
		})},
		Input: []Token{
			{Value: "a", Type: Text},
		},
		Output: []Token{},
		Error:  io.EOF,
	}, {
		Name:    "MergeTokensFilter",
		Filters: Filters{MergeTokensFilter},
		Input: []Token{
			{Value: "a", Type: Text},
			{Value: "b", Type: Text},
			{Value: "c", Type: Text},
		},
		Output: []Token{
			{Value: "abc", Type: Text},
		},
	}, {
		Name:    "RemoveEmptiesFilter",
		Filters: Filters{RemoveEmptiesFilter},
		Input: []Token{
			{Value: "a", Type: Text},
			{Value: "", Type: Text},
			{Value: "c", Type: Text},
		},
		Output: []Token{
			{Value: "a", Type: Text},
			{Value: "c", Type: Text},
		},
	}, {
		Name:    "MergeTokensFilter -> RemoveEmptiesFilter",
		Filters: Filters{RemoveEmptiesFilter, MergeTokensFilter},
		Input: []Token{
			{Value: "a", Type: Text},
			{Value: "", Type: Text},
			{Value: "c", Type: Text},
		},
		Output: []Token{
			{Value: "ac", Type: Text},
		},
	}} {
		err := testFilters(t, item.Filters, item.Input, item.Output, item.Name)
		assert.Equal(t, item.Error, err, item.Name)
	}
}

func testFilters(t *testing.T, filters Filters, input, expected []Token,
	name string) error {

	pos := 0

	filter := filters.Filter(func(token Token) error {
		assert.Equal(t, expected[pos], token, name)
		pos++
		return nil
	})

	for _, inToken := range input {
		err := filter(inToken)
		if err != nil {
			return err
		}
	}

	return nil
}
