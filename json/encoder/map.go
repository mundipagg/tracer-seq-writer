package encoder

import (
	"github.com/json-iterator/go"
	"unsafe"
)

type Map struct {
	Strategy func(string) string
}

func (enc *Map) IsEmpty(ptr unsafe.Pointer) bool {
	s := (*string)(ptr)
	return s == nil || *s == ""
}

func (enc *Map) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if enc.IsEmpty(ptr) {
		stream.WriteString("")
	} else {
		s := (*string)(ptr)
		stream.WriteString(enc.Strategy(*s))
	}
}
