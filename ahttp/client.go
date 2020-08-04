package ahttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// the parameter "query" should be nil or a struct,
// "out" should be a &string, or a &struct.
func Get(path string, query interface{}, out interface{}) error {
	var resp, err = requestGet(path, query)
	if err != nil {
		return err
	}

	return handleGet(resp, out)
}

func requestGet(path string, query interface{}) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("path is empty")
	}

	var url strings.Builder
	url.WriteString(path)

	// "query" can be ignored.
	if query != nil {
		str, err := MarshalQuery(query)
		if err != nil {
			return "", err
		}
		if len(str) > 0 {
			url.WriteString("?")
			url.WriteString(str)
		}
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}

	// NOTE: the status code isn't 200, regard as an error.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("returned status code is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func handleGet(resp string, out interface{}) error {
	if out == nil {
		return fmt.Errorf("output parameter is nil")
	}

	var ptr = reflect.ValueOf(out)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("output parameter isn't a ptr")
	}

	var kind = ptr.Elem().Kind()
	if kind == reflect.Struct {
		//
		return nil

	} else if kind == reflect.String {
		*(out.(*string)) = resp
		return nil

	} else {
		return fmt.Errorf("unsupported output type")
	}
}
