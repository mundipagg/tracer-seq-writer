// +build unit

package encoder

import (
	"bytes"
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestMap_Encode(t *testing.T) {
	t.Parallel()
	t.Run("when the input is 'empty'", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		expected := ""
		subject := &Map{
			Strategy: func(s string) string {
				called++
				return expected
			},
		}
		input := ""
		pointer := reflect.ValueOf(&input).Pointer()
		buf := &bytes.Buffer{}
		stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
		subject.Encode(unsafe.Pointer(pointer), stream)
		_ = stream.Flush()
		is.Equal(strconv.Quote(expected), buf.String())
	})
	t.Run("when the input is not 'empty'", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		called := 0
		expected := "output"
		subject := &Map{
			Strategy: func(s string) string {
				called++
				return expected
			},
		}
		input := "input"
		pointer := reflect.ValueOf(&input).Pointer()
		buf := &bytes.Buffer{}
		stream := jsoniter.NewStream(jsoniter.ConfigFastest, buf, 100)
		subject.Encode(unsafe.Pointer(pointer), stream)
		_ = stream.Flush()
		is.Equal(strconv.Quote(expected), buf.String())
	})
}

func TestMap_IsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	subject := &Map{}
	input := "input"
	pointer := reflect.ValueOf(&input).Pointer()
	is.False(subject.IsEmpty(unsafe.Pointer(pointer)))
	input = ""
	pointer = reflect.ValueOf(&input).Pointer()
	is.True(subject.IsEmpty(unsafe.Pointer(pointer)))
	is.True(subject.IsEmpty(nil))
}
