package ahttp

import (
	"src/clog"
)

func DefHandler(trans *HTTPTrans) {

	switch trans.ReqPath {
	case "/alpha":
		clog.I("call func 'alpha'")
		handleAlpha(trans)

	case "/beta":
		clog.I("call func 'beta'")
		handleBeta(trans)

	default:
		clog.E("unknown func '%s'", trans.ReqPath)
	}
}

func handleAlpha(trans *HTTPTrans) {

	var query AlphaQuery
	var resp AlphaResp

	//
	var err = UnmarshalQuery(trans.ReqQuery, &query)
	if err != nil {
		resp.ErrCode = AlphaCodeParamErr
		resp.ErrDesc = err.Error()

		trans.RespBody = resp

		return
	}

	//
	resp.ErrCode = AlphaCodeOkay
	resp.ErrDesc = ""
	resp.Int = query.Int
	resp.Str = query.Str

	trans.RespBody = resp
}

func handleBeta(trans *HTTPTrans) {
}
