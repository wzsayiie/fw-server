package api

import (
	"src/clog"
)

func DefHTTPHandler(trans *HTTPTrans) {

	switch trans.ReqPath {
	case "/alpha":
		clog.I("call api 'alpha'")
		handleAlpha(trans)

	case "/beta":
		clog.I("call api 'beta'")
		handleBeta(trans)

	default:
		clog.E("unknown api '%s'", trans.ReqPath)
	}
}

func handleAlpha(trans *HTTPTrans) {

	var request AlphaRequest
	var response AlphaResponse

	//
	var err = UnmarshalQuery(trans.ReqQuery, &request)
	if err != nil {
		response.ErrCode = AlphaCodeQueryErr
		response.ErrDesc = err.Error()
		trans.RespBody = response
		return
	}

	//
	response.ErrCode = AlphaCodeOkay
	response.ErrDesc = ""
	response.Int = request.Int
	response.Str = request.Str

	trans.RespBody = response
}

func handleBeta(trans *HTTPTrans) {
}
