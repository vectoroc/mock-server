package matcher

import (
	"mock-server/model"
	"net/http"
	"strings"
)

// NormalizeRequest makes it possible to match GET and get requests, Content-Type and content-type headers and so on
func NormalizeRequest(req *model.HttpRequest) {
	replace := make(map[string]string)
	for k := range req.Headers {
		canonical := http.CanonicalHeaderKey(k)
		if k != canonical {
			replace[k] = canonical
		}
	}

	for k, canon := range replace {
		req.Headers[canon] = req.Headers[k]
		delete(req.Headers, k)
	}

	req.Method = strings.ToUpper(req.Method)

	if req.Path > "" && req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
}
