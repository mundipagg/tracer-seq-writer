// +build unit

package encoder

import (
	"bytes"
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestStruct_IsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	subject := &Struct{}
	input := struct{}{}
	pointer := reflect.ValueOf(&input).Pointer()
	is.False(subject.IsEmpty(unsafe.Pointer(pointer)))
}

func TestStruct_Encode(t *testing.T) {
	t.Parallel()
	t.Run("when the field does not have a json tag", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		input := struct {
			A int
		}{
			A: 15,
		}
		called := 0
		subject := &Struct{
			Strategy: func(s string) string {
				called++
				return strings.ToLower(s)
			},
			Type: reflect.TypeOf(input),
		}
		buf := &bytes.Buffer{}
		stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
		subject.Encode(unsafe.Pointer(reflect.ValueOf(&input).Pointer()), stream)
		stream.Flush()
		is.Equal(`{"a":15}`, buf.String(), "it should change the name of the field")
		is.Equal(1, called, "it should call the strategy exactly one time")
	})
	t.Run("when the field does not have omitempty on it's tag", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		input := struct {
			A int `json:"SuperParameter"`
		}{
			A: 15,
		}
		called := 0
		subject := &Struct{
			Strategy: func(s string) string {
				called++
				return strings.ToUpper(s)
			},
			Type: reflect.TypeOf(input),
		}
		buf := &bytes.Buffer{}
		stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
		subject.Encode(unsafe.Pointer(reflect.ValueOf(&input).Pointer()), stream)
		stream.Flush()
		is.Equal(`{"SUPERPARAMETER":15}`, buf.String(), "it should change the name of the field")
		is.Equal(1, called, "it should call the strategy exactly one time")
	})
	t.Run("when the field does have omitempty on it's tag", func(t *testing.T) {
		t.Parallel()
		t.Run("but the field is not empty", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			input := struct {
				A int `json:"SuperParameter,omitempty"`
				B int `json:"SuperParameterB"`
			}{
				A: 15,
			}
			called := 0
			subject := &Struct{
				Strategy: func(s string) string {
					called++
					return strings.ToUpper(s)
				},
				Type: reflect.TypeOf(input),
			}
			buf := &bytes.Buffer{}
			stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
			subject.Encode(unsafe.Pointer(reflect.ValueOf(&input).Pointer()), stream)
			stream.Flush()
			is.Equal(`{"SUPERPARAMETER":15,"SUPERPARAMETERB":0}`, buf.String(), "it should change the name of the field")
			is.Equal(2, called, "it should call the strategy exactly two times")
		})
		t.Run("and the field is empty", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			input := struct {
				A *int `json:"SuperParameter,omitempty"`
			}{
				A: nil,
			}
			called := 0
			subject := &Struct{
				Strategy: func(s string) string {
					called++
					return strings.ToUpper(s)
				},
				Type: reflect.TypeOf(input),
			}
			buf := &bytes.Buffer{}
			stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
			subject.Encode(unsafe.Pointer(reflect.ValueOf(&input).Pointer()), stream)
			stream.Flush()
			is.Equal(`{}`, buf.String(), "it should change the name of the field")
			is.Equal(0, called, "it should not call the strategy")
		})
	})
}
