// +build unit

package seq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry_Add(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	subject := Entry{
		"A": 15,
	}
	subject.Add("B", 16)
	expected := Entry{
		"A": 15,
		"B": 16,
	}
	is.Equal(expected, subject, "it should update the entry")
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("when an interface slice is given", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		type S struct {
			B int
		}
		items := []interface{}{
			"A",
			Entry{
				"X": "Y",
			},
			nil,
			S{15},
			&S{16},
		}
		actual := NewEntry(items)
		expected := Entry{
			"string": "A",
			"X":      "Y",
			"S":      S{15},
			"S1":     S{16},
		}
		is.Equal(expected, actual, "it should return the expected entry")
	})
	t.Run("when multiple items are given", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		type S struct {
			B int
		}
		items := []interface{}{
			"A",
			Entry{
				"X": "Y",
			},
			S{15},
			&S{16},
		}
		actual := NewEntry(items...)
		expected := Entry{
			"string": "A",
			"X":      "Y",
			"S":      S{15},
			"S1":     S{16},
		}
		is.Equal(expected, actual, "it should return the expected entry")
	})
}

func TestMerge(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	inputA := Entry{
		"A": 15,
		"B": 16,
		"C": 17,
	}
	inputB := Entry{
		"A": 20,
		"E": 21,
		"F": 22,
	}
	actual := Merge(inputA, inputB)
	expected := Entry{
		"A":  15,
		"B":  16,
		"C":  17,
		"A1": 20,
		"E":  21,
		"F":  22,
	}
	is.Equal(expected, actual, "it should return the expected value")
}
