package api

import (
	"encoding/xml"
)

type AlphaReqQuery struct {
	Int int64  `query:"int"`
	Str string `query:"str"`
}

const (
	AlphaRespEOkay     = 0
	AlphaRespEQueryErr = 1
	/* ... */
)

type AlphaRespBody struct {
	XMLName xml.Name `json:"-" xml:"root"`

	ErrCode int64  `json:"errcode" xml:"errcode"`
	ErrDesc string `json:"errdesc" xml:"errdesc"`

	Int int64  `json:"int" xml:"int"`
	Str string `json:"str" xml:"str"`
}
