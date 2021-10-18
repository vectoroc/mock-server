package server

import (
	"mock-server/model"
	"net/http"
)

func httpCookieToKeyToValue(cookies []*http.Cookie) model.KeyToValue {
	res := model.KeyToValue{Values: make(map[string]string, len(cookies))}
	for _, c := range cookies {
		res.Values[c.Name] = c.Value
	}
	return res
}

func httpHeadersToKeyToMultiValue(req *http.Request) model.KeyToMultiValue {
	// http.Request does not contain Host header ? TODO: test it
	res := model.KeyToMultiValue{Values: req.Header.Clone()}
	res.Values["Host"] = []string{req.Host}
	return res
}
