// +build unit

package json

import (
	"github.com/modern-go/reflect2"
	"github.com/mundipagg/tracer-seq-writer/json/encoder"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCaseStrategyExtension_CreateMapKeyEncoder(t *testing.T) {
	t.Parallel()
	t.Run("when the key is a string", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		subject := &CaseStrategyExtension{
			Strategy: func(s string) string {
				called++
				return s
			},
		}
		typ := reflect2.Type2(reflect.TypeOf(""))
		actual := subject.CreateMapKeyEncoder(typ)
		is.IsType(&encoder.Map{}, actual, "it should return a Map encoder")
		is.Equal(0, called, "it should not call the strategy")
	})
	t.Run("when the key is not a string", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		subject := &CaseStrategyExtension{
			Strategy: func(s string) string {
				called++
				return s
			},
		}
		typ := reflect2.Type2(reflect.TypeOf(1))
		actual := subject.CreateMapKeyEncoder(typ)
		is.Nil(actual, "it should return nil")
		is.Equal(0, called, "it should not call the strategy")
	})
}

func TestCaseStrategyExtension_CreateEncoder(t *testing.T) {
	t.Parallel()
	t.Run("when the key is a struct", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		subject := &CaseStrategyExtension{
			Strategy: func(s string) string {
				called++
				return s
			},
		}
		typ := reflect2.Type2(reflect.TypeOf(struct {
		}{}))
		actual := subject.CreateEncoder(typ)
		is.IsType(&encoder.Struct{}, actual, "it should return a Struct encoder")
		is.Equal(0, called, "it should not call the strategy")
	})
	t.Run("when the key is not a struct", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		subject := &CaseStrategyExtension{
			Strategy: func(s string) string {
				called++
				return s
			},
		}
		typ := reflect2.Type2(reflect.TypeOf(1))
		actual := subject.CreateEncoder(typ)
		is.Nil(actual, "it should return nil")
		is.Equal(0, called, "it should not call the strategy")
	})
}
