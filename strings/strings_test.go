// +build unit

package strings

import (
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBlank(t *testing.T) {
	t.Parallel()
	t.Run("when string is empty", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		result := IsBlank("")
		is.True(result, "it should return true")
	})
	t.Run("when string only contains empty characters", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		result := IsBlank("     ")
		is.True(result, "it should return true")
	})
	t.Run("when string contain non whitespace characters", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		result := IsBlank(fake.CharactersN(6))
		is.False(result, "it should return false")
	})
}

func TestToPascalCase(t *testing.T) {
	t.Parallel()
	t.Run("when a snake case string is passed", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		in := "pascal_case_word"
		expected := "PascalCaseWord"
		is.Equal(expected, ToPascalCase(in), "should return the word as pascal")
	})
}

func TestCapitalize(t *testing.T) {
	t.Parallel()
	t.Run("when the string has more than one character", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		is.Equal("DeCaPiTaLiZe", Capitalize("deCaPiTaLiZe"))
	})
	t.Run("when the string has one character", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		is.Equal("D", Capitalize("d"))
	})
}
