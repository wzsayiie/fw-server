package ahttp

import (
	"encoding/xml"
)

type AlphaQuery struct {
	Int int64  `query:"int"`
	Str string `query:"str"`
}

const (
	AlphaCodeOkay     = 0
	AlphaCodeParamErr = 1
	/* ... */
)

type AlphaResp struct {
	XMLName xml.Name `json:"-" xml:"root"`

	ErrCode int64  `json:"errcode" xml:"errcode"`
	ErrDesc string `json:"errdesc" xml:"errdesc"`

	Int int64  `json:"int" xml:"int"`
	Str string `json:"str" xml:"str"`
}
