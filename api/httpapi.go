package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"src/log"
)

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
		handle(handler, resp, req)
	})

	// don't specify a clear ip,
	// so that the local loopback address can be used during the test.
	// NOTE:
	// if the host has multiple network adapters, can't specify which adapter to use.
	var addr = fmt.Sprintf(":%d", port)

	var err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.E("%s", err.Error())
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

func handle(handler func(trans *HTTPTrans), resp http.ResponseWriter, req *http.Request) {

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
	// NOTE: json is default.
	var val = fmt.Sprintf(`{"errcode":%d,"errdesc":"%s"}`, code, desc)

	log.E("resp {")
	log.E("  from: %s", req.RemoteAddr)
	log.E("  path: %s", req.URL.String())
	log.E("  resp: %s", val)
	log.E("}")

	resp.Write([]byte(val))
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
		if len(v) > 1 {
			return nil, fmt.Errorf("duplicate paramter")
		}

		ret[k] = v[0]
	}

	return ret, nil
}
