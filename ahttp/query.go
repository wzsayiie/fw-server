package ahttp

import (
	"fmt"
	"reflect"
	"strconv"
)

// unmarshal query map into a struct.
// only type bool, int, uint, float and string are supported.

func UnmarshalQuery(query map[string]string, out interface{}) error {
	if query == nil || out == nil {
		return fmt.Errorf("parameter is nil")
	}

	var ptrVal = reflect.ValueOf(out)
	if ptrVal.Kind() != reflect.Ptr {
		return fmt.Errorf("output parameter isn't a ptr")
	}

	var stcVal = ptrVal.Elem()
	if stcVal.Kind() != reflect.Struct {
		return fmt.Errorf("output parameter isn't point to a struct")
	}

	for i := 0; i < stcVal.NumField(); i++ {
		var fieldVal = stcVal.Field(i)
		var fieldTag = fieldTag(stcVal, i)
		if fieldTag == "-" {
			// ignore this field.
			continue
		}

		var str, exist = query[fieldTag]
		if !exist {
			// if there is no corresponding new value, keep the origin value.
			continue
		}

		var err = setField(&fieldVal, fieldTag, str)
		if err != nil {
			return err
		}
	}

	return nil
}

func fieldTag(stc reflect.Value, idx int) string {
	var typ = stc.Type().Field(idx)
	var tag = typ.Tag.Get("query")

	if tag == "" {
		return typ.Name
	} else {
		return tag
	}
}

func setField(field *reflect.Value, name string, str string) error {
	switch field.Kind() {
	case reflect.Bool:
		var val, err = strconv.ParseBool(str)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type bool", str, name)
		}

		field.SetBool(val)

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		var val, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type int", str, name)
		}

		field.SetInt(val)

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		var val, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type uint", str, name)
		}

		field.SetUint(val)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		var val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type float", str, name)
		}

		field.SetFloat(val)

	case reflect.String:
		field.SetString(str)

	default:
		return fmt.Errorf("can't set '%s' cause unsupported type", name)
	}

	return nil
}
