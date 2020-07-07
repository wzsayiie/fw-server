package api

import (
	"src/log"
)

func DefHTTPHandler(trans *HTTPTrans) {

	switch trans.ReqPath {
	case "/alpha":
		log.I("call api 'alpha'")
		handleAlpha(trans)

	case "/beta":
		log.I("call api 'beta'")
		handleBeta(trans)

	default:
		log.E("unknown api '%s'", trans.ReqPath)
	}
}

func handleAlpha(trans *HTTPTrans) {

	var reqQuery AlphaReqQuery
	var respBody AlphaRespBody

	var err = UnmarshalQuery(trans.ReqQuery, &reqQuery)
	if err == nil {
		respBody.ErrCode = AlphaRespEOkay
		respBody.ErrDesc = ""
		respBody.Int = reqQuery.Int
		respBody.Str = reqQuery.Str

	} else {
		respBody.ErrCode = AlphaRespEQueryErr
		respBody.ErrDesc = err.Error()
	}

	trans.RespBody = respBody
}

func handleBeta(trans *HTTPTrans) {
}
