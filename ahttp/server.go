package ahttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"src/clog"
	"src/local"
)

type HTTPTrans struct {
	FromAddr string
	ReqPath  string
	ReqQuery map[string]string

	// user program need to assign this field as response body.
	RespBody interface{}
}

func Serve(port uint16, handler func(trans *HTTPTrans)) {

	if handler == nil {
		clog.E("http handler is nil")
		return
	}

	// print local address information.
	// NOTE:
	// only the first network adpater found is printed.
	// by default, there is only one adapter on this host.
	var ip4 string
	var ip6 string
	local.HostIPs(&ip4, &ip6)

	clog.I("http ready on {")
	clog.I("  [%s]:%d", ip6, port)
	clog.I("  %s:%d", ip4, port)
	clog.I("}")

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		handle(handler, resp, req)
	})

	// don't specify a clear ip,
	// so that the local loopback address can be used during the test.
	// NOTE:
	// if the host has multiple network adapters, can't specify which adapter to use.
	var addr = fmt.Sprintf(":%d", port)

	var err = http.ListenAndServe(addr, nil)
	if err != nil {
		clog.E("%s", err)
	}
}

const (
	CodeInvalidMethod = -1
	CodeQueryError    = -2
	CodeNoResponse    = -3
)

func handle(handler func(trans *HTTPTrans), resp http.ResponseWriter, req *http.Request) {

	logReq(req)

	// if user use a browser for testing, the browser will request the icon by this path.
	if req.URL.Path == "/favicon.ico" {
		return
	}

	if req.Method != "GET" {
		// only method "GET" can be supported.
		logRespExcept(req, resp, CodeInvalidMethod, "only 'Get' can be supported")
		return
	}

	var query, err = filterQuery(req.URL.Query())
	if err != nil {
		logRespExcept(req, resp, CodeQueryError, err.Error())
		return
	}

	var trans = HTTPTrans{
		FromAddr: req.RemoteAddr,
		ReqPath:  req.URL.Path,
		ReqQuery: query,
	}
	handler(&trans)

	if trans.RespBody == nil {
		logRespExcept(req, resp, CodeNoResponse, "the program didn't respond")
		return
	}

	logRespNormal(req, resp, trans.RespBody)
}

func logReq(req *http.Request) {
	clog.I("REQ {")
	clog.I("  FROM: %s", req.RemoteAddr)
	clog.I("  MTHD: %s", req.Method)
	clog.I("  PATH: %s", req.URL.String())
	clog.I("}")
}

func logRespExcept(req *http.Request, resp http.ResponseWriter, code int64, desc string) {

	type Body struct {
		XMLName xml.Name `json:"-" xml:"root"`

		ErrCode int64  `json:"errcode" xml:"errcode"`
		ErrDesc string `json:"errdesc" xml:"errdesc"`
	}

	var body = Body{
		ErrCode: code,
		ErrDesc: desc,
	}
	logRespNormal(req, resp, body)
}

func logRespNormal(req *http.Request, resp http.ResponseWriter, body interface{}) {
	// NOTE: json is default.
	var val, _ = json.Marshal(body)

	clog.I("resp {")
	clog.I("  from: %s", req.RemoteAddr)
	clog.I("  path: %s", req.URL.String())
	clog.I("  resp: %s", val)
	clog.I("}")

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
