package api

import (
	"encoding/xml"
)

type AlphaRequest struct {
	Int int64  `query:"int"`
	Str string `query:"str"`
}

const (
	AlphaCodeOkay     = 0
	AlphaCodeQueryErr = 1
	/* ... */
)

type AlphaResponse struct {
	XMLName xml.Name `json:"-" xml:"root"`

	ErrCode int64  `json:"errcode" xml:"errcode"`
	ErrDesc string `json:"errdesc" xml:"errdesc"`

	Int int64  `json:"int" xml:"int"`
	Str string `json:"str" xml:"str"`
}
