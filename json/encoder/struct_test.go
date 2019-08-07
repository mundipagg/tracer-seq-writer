// +build unit

package encoder

import (
	"bytes"
	"errors"
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

type V struct {
	A int
}

func (V) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal("custom")
}

func TestStruct_Encode(t *testing.T) {
	t.Parallel()
	t.Run("when the value implements the json.Marshaller interface", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		input := V{15}
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
		is.Equal(`"custom"`, buf.String(), "it should change the name of the field")
		is.Equal(0, called, "it should not call the strategy ")
	})
	t.Run("when the value implements the error interface", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		input := errors.New("error")
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
		ptr := reflect.New(reflect.TypeOf(input))
		ptr.Elem().Set(reflect.ValueOf(input))
		subject.Encode(unsafe.Pointer(ptr.Pointer()), stream)
		stream.Flush()
		is.Equal(`"error"`, buf.String(), "it should change the name of the field")
		is.Equal(0, called, "it should not call the strategy ")
	})
	t.Run("when the field does not have a json tag", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		input := struct {
			A int
			b int
		}{
			A: 15,
			b: 15,
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
			v := 16
			input := struct {
				A int  `json:"SuperParameterA"`
				B *int `json:"SuperParameterB,omitempty"`
				C int `json:"SuperParameterC,omitempty"`
				D *int `json:"SuperParameterD,omitempty"`
				E int  `json:"SuperParameterE"`
			}{
				A: 15,
				D: &v,
				E: 17,
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
			is.Equal(`{"SUPERPARAMETERA":15,"SUPERPARAMETERD":16,"SUPERPARAMETERE":17}`, buf.String(), "it should change the name of the field")
			is.Equal(3, called, "it should call the strategy exactly three times")
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
