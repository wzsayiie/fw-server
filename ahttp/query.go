package ahttp

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// NOTE:
// only type BOOL, INT, UINT, FLOAT and STRING are supported.

//NOTE: the parameter "out" should be a &struct.
func UnmarshalQuery(query map[string]string, out interface{}) error {
	if query == nil || out == nil {
		return fmt.Errorf("parameter is nil")
	}

	var ptr = reflect.ValueOf(out)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("output parameter isn't a ptr")
	}

	var stc = ptr.Elem()
	if stc.Kind() != reflect.Struct {
		return fmt.Errorf("output parameter isn't point to a struct")
	}

	for i := 0; i < stc.NumField(); i++ {

		var tag = fieldTag(stc, i)
		if tag == "-" {
			// ignore this field.
			continue
		}

		var str, inc = query[tag]
		if !inc {
			// if there is no corresponding new value, keep the origin value.
			continue
		}

		var mem = stc.Field(i)
		var err = setField(&mem, tag, str)
		if err != nil {
			return err
		}
	}

	return nil
}

//NOTE: the parameter "dat" should be a struct.
func MarshalQuery(dat interface{}) (string, error) {
	if dat == nil {
		return "", fmt.Errorf("input parameter is nil")
	}

	var stc = reflect.ValueOf(dat)
	if stc.Kind() != reflect.Struct {
		return "", fmt.Errorf("input parameter is not a struct")
	}

	var buf strings.Builder

	for i := 0; i < stc.NumField(); i++ {
		var tag = fieldTag(stc, i)
		if tag == "-" {
			// ignore this field.
			continue
		}

		var val = stc.Field(i)
		var str = strVal(val)

		if buf.Len() > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(url.QueryEscape(tag))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(str))
	}

	return buf.String(), nil
}

func fieldTag(stc reflect.Value, idx int) string {
	var mem = stc.Type().Field(idx)
	var tag = mem.Tag.Get("query")

	if tag == "" {
		// if the tag didn't specified, return the field name as tag.
		return mem.Name
	} else {
		return tag
	}
}

func setField(field *reflect.Value, name string, str string) error {
	if isBool(*field) {
		var val, err = strconv.ParseBool(str)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type bool", str, name)
		}

		field.SetBool(val)

	} else if isInt(*field) {
		var val, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type int", str, name)
		}

		field.SetInt(val)

	} else if isUint(*field) {
		var val, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type uint", str, name)
		}

		field.SetUint(val)

	} else if isFloat(*field) {
		var val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("can't unmarshal '%s' into '%s' of type float", str, name)
		}

		field.SetFloat(val)

	} else if isString(*field) {
		// no query escape decode operation here,
		// because the operation was done by http.Request .
		field.SetString(str)

	} else {
		return fmt.Errorf("can't set '%s' cause unsupported type", name)
	}

	return nil
}

func strVal(val reflect.Value) string {
	if isBool(val) {
		return strconv.FormatBool(val.Bool())

	} else if isInt(val) {
		return strconv.FormatInt(val.Int(), 10)

	} else if isUint(val) {
		return strconv.FormatUint(val.Uint(), 10)

	} else if isFloat32(val) {
		return strconv.FormatFloat(val.Float(), 'f', -1, 32)

	} else if isFloat64(val) {
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)

	} else if isString(val) {
		return val.String()

	} else {
		return ""
	}
}

func isBool(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Bool:
		return true

	default:
		return false
	}
}

func isInt(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return true

	default:
		return false
	}
}

func isUint(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return true

	default:
		return false
	}
}

func isFloat(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return true

	default:
		return false
	}
}

func isFloat32(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Float32:
		return true

	default:
		return false
	}
}

func isFloat64(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Float64:
		return true

	default:
		return false
	}
}

func isString(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.String:
		return true

	default:
		return false
	}
}
