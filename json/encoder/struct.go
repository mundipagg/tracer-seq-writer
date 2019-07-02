package encoder

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
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
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
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
		for i := 0; i < numFields-1; i++ {
			fv := v.Field(i)
			ft := changer.Type.Field(i)
			if changer.writeField(ft, fv, stream) {
				stream.WriteMore()
			}
		}
		fv := v.Field(numFields - 1)
		ft := changer.Type.Field(numFields - 1)
		changer.writeField(ft, fv, stream)
		stream.WriteObjectEnd()
	}
}

func (changer *Struct) writeField(structField reflect.StructField, value reflect.Value, stream *jsoniter.Stream) bool {
	if !value.CanInterface() {
		return false
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
