package api

import (
	"encoding/xml"
)

type XXXReqQuery struct {
	Int int    `query:"int"`
	Str string `query:"str"`
}

const (
	XXXRespEOkay     = 0
	XXXRespEQueryErr = 1
	/* ... */
)

type XXXRespBody struct {
	XMLName xml.Name `json:"-" xml:"root"`

	ErrCode int    `json:"errcode" xml:"errcode"`
	ErrDesc string `json:"errdesc" xml:"errdesc"`

	Int int    `json:"int" xml:"int"`
	Str string `json:"str" xml:"str"`
}
