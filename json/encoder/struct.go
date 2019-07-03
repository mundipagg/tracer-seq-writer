package encoder

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

type Struct struct {
	Type     reflect.Type
	Strategy func(string) string
}

func (changer *Struct) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (changer *Struct) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	beforeBuffer := stream.Buffer()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "a error occurred while serialization of '%v', error: '%v'", changer.Type.Name(), err)
			stream.SetBuffer(beforeBuffer)
		}
	}()
	v := reflect.NewAt(changer.Type, ptr).Elem()
	switch value := v.Interface().(type) {
	case json.Marshaler:
		valueJ, _ := value.MarshalJSON()
		var valueM interface{}
		_ = json.Unmarshal(valueJ, &valueM)
		stream.WriteVal(valueM)
	case error:
		stream.WriteString(value.Error())
	default:
		stream.WriteObjectStart()
		numFields := v.NumField()
		if numFields > 0 {
			fv := v.Field(0)
			ft := changer.Type.Field(0)
			first := !changer.writeField(ft, fv, stream, true)
			for i := 1; i < numFields; i++ {
				fv := v.Field(i)
				ft := changer.Type.Field(i)
				if changer.writeField(ft, fv, stream, first) {
					first = false
				}
			}
		}
		stream.WriteObjectEnd()
	}
}

func (changer *Struct) writeField(structField reflect.StructField, value reflect.Value, stream *jsoniter.Stream, first bool) bool {
	if !value.CanInterface() {
		return false
	}
	if !first {
		stream.WriteMore()
	}
	tag := strings.TrimSpace(structField.Tag.Get("json"))
	if len(tag) == 0 {
		stream.WriteObjectField(changer.Strategy(structField.Name))
		stream.WriteVal(value.Interface())
	} else {
		pieces := strings.Split(tag, ",")
		if len(pieces) > 1 {
			if pieces[1] == "omitempty" {
				isNil := func() (isNil bool) {
					defer func() {
						if recover() != nil {
							isNil = false
						}
					}()
					isNil = value.IsNil()
					return isNil
				}()
				if isNil {
					return false
				}
			}
		}
		stream.WriteObjectField(changer.Strategy(pieces[0]))
		stream.WriteVal(value.Interface())
	}
	return true
}
