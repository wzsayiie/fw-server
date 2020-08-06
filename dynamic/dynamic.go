package dynamic

import (
	"reflect"
)

type Kind uint

const (
	Other Kind = iota

	Nil
	StringPtr
	StructPtr
)

func KindOf(a interface{}) Kind {

	var val = reflect.ValueOf(a)

	switch val.Kind() {
	case reflect.Invalid:
		return Nil

	case reflect.Ptr:
		return checkPtrKind(val.Elem())

	default:
		return Other
	}
}

func checkPtrKind(val reflect.Value) Kind {

	switch val.Kind() {
	case reflect.String:
		return StringPtr

	case reflect.Struct:
		return StructPtr

	default:
		return Other
	}
}
