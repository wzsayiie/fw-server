package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"src/log"
	"strconv"
)

// unmarshal query map into a struct:
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
		var fieldTag = structFieldTag(stcVal, i)

		var str, exist = query[fieldTag]
		if !exist {
			// if there is no corresponding new value, keep the origin value.
			continue
		}

		var err = setStructField(&fieldVal, fieldTag, str)
		if err != nil {
			return err
		}
	}

	return nil
}

func structFieldTag(stc reflect.Value, idx int) string {
	var typ = stc.Type().Field(idx)
	var tag = typ.Tag.Get("query")

	if tag == "" {
		return typ.Name
	} else {
		return tag
	}
}

func setStructField(field *reflect.Value, name string, str string) error {
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

// a http server:

type HTTPTrans struct {
	FromAddr string
	ReqPath  string
	ReqQuery map[string]string

	// user program need to assign this field as response body.
	RespBody interface{}
}

func HTTPServe(port uint16, handler func(trans *HTTPTrans)) {

	if handler == nil {
		log.E("http handler is nil")
		return
	}

	// print local address infomation.
	// NOTE:
	// only the first network adpater found is printed.
	// by default, there is only one adapter on this host.
	log.I("http ready on {")
	for _, v := range favLocalAddrs() {
		log.I("  %s :%d", v.String(), port)
	}
	log.I("}")

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		httpHandle(handler, resp, req)
	})

	// don't specify a clear ip,
	// so that the local loopback address can be used during the test.
	// NOTE:
	// if the host has multiple network adapters, can't specify which adapter to use.
	var addr = fmt.Sprintf(":%d", port)

	var err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.E("%s", err)
	}
}

func favLocalAddrs() []net.Addr {

	var intfs, intfErr = net.Interfaces()
	if intfErr != nil {
		return nil
	}

	for _, v := range intfs {
		if v.Flags&net.FlagUp == 0 /* isn't working */ {
			continue
		}
		if v.Flags&net.FlagLoopback != 0 /* is loopback */ {
			continue
		}

		var addrs, addrErr = v.Addrs()
		if addrErr != nil {
			continue
		}

		// there are usually two valus, ipv4 and ipv6.
		if len(addrs) != 0 {
			return addrs
		}
	}
	return nil
}

func httpHandle(handler func(trans *HTTPTrans), resp http.ResponseWriter, req *http.Request) {

	logReq(req)

	// if user use a browser for testing, the browser will request the icon by this path.
	if req.URL.Path == "/favicon.ico" {
		return
	}

	if req.Method != "GET" {
		// only method "GET" can be supported.
		logRespExcept(req, resp, -1, "only method 'Get' can be supported")
		return
	}

	var query, err = filterQuery(req.URL.Query())
	if err != nil {
		logRespExcept(req, resp, -2, err.Error())
		return
	}

	var trans = HTTPTrans{
		FromAddr: req.RemoteAddr,
		ReqPath:  req.URL.Path,
		ReqQuery: query,
	}
	handler(&trans)

	if trans.RespBody == nil {
		logRespExcept(req, resp, -3, "the program didn't return the result")
		return
	}

	logRespNormal(req, resp, trans.RespBody)
}

func logReq(req *http.Request) {
	log.I("REQ {")
	log.I("  FROM: %s", req.RemoteAddr)
	log.I("  MTHD: %s", req.Method)
	log.I("  PATH: %s", req.URL.String())
	log.I("}")
}

func logRespExcept(req *http.Request, resp http.ResponseWriter, code int, desc string) {

	type Body struct {
		XMLName xml.Name `json:"-" xml:"root"`

		ErrCode int64  `json:"errcode" xml:"errcode"`
		ErrDesc string `json:"errdesc" xml:"errdesc"`
	}

	var body = Body{
		ErrCode: int64(code),
		ErrDesc: desc,
	}
	logRespNormal(req, resp, body)
}

func logRespNormal(req *http.Request, resp http.ResponseWriter, body interface{}) {
	// NOTE: json is default.
	var val, _ = json.Marshal(body)

	log.I("resp {")
	log.I("  from: %s", req.RemoteAddr)
	log.I("  path: %s", req.URL.String())
	log.I("  resp: %s", val)
	log.I("}")

	resp.Write([]byte(val))
}

func filterQuery(raw url.Values) (map[string]string, error) {
	var ret = make(map[string]string)

	for k, v := range raw {
		// parameters with the same name are not supported.
		if len(v) > 1 {
			return nil, fmt.Errorf("duplicate paramter")
		}

		ret[k] = v[0]
	}

	return ret, nil
}
