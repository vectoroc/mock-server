package matcher

import (
	"mock-server/model"
	"net/http"
	"strings"
)

func NormalizeRequest(req *model.HttpRequest) {
	replace := make(map[string]string)
	for k := range req.Headers.Values {
		canonical := http.CanonicalHeaderKey(k)
		if k != canonical {
			replace[k] = canonical
		}
	}

	for k, canon := range replace {
		req.Headers.Values[canon] = req.Headers.Values[k]
		delete(req.Headers.Values, k)
	}

	req.Method = strings.ToUpper(req.Method)

	if req.Path > "" && req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
}
