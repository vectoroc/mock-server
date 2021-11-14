package server

import (
	"mock-server/model"
	"net/http"
)

func httpCookieToKeyToValue(cookies []*http.Cookie) model.KeyToValue {
	res := model.KeyToValue{}
	for _, c := range cookies {
		res[c.Name] = c.Value
	}
	return res
}

func httpHeadersToKeyToMultiValue(req *http.Request) model.KeyToMultiValue {
	// http.Request does not contain Host header ? TODO: test it
	res := model.KeyToMultiValue(req.Header.Clone())
	res["Host"] = []string{req.Host}
	return res
}
