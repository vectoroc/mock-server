package server

import (
	"io"
	"mock-server/matcher"
	"mock-server/model"
	"net/http"
	"net/url"
	"strings"
)

func matchHttpRequest(exp *model.HttpRequest, req *http.Request) (bool, error) {

	if exp.Method != "" && req.Method != strings.ToUpper(exp.Method) {
		return false, nil
	}

	if exp.Body != nil {
		contentType := req.Header.Get("Content-Type")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return false, err
		}

		if !MatchBody(string(body), contentType, exp.Body) {
			return false, err
		}
	}

	if exp.Path != "" && exp.Path != req.URL.Path {
		return false, nil
	}

	for name, v := range exp.QueryStringParameters.Values {
		queryParams, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			return false, err
		}

		if !matcher.Equal(v, queryParams[name]) {
			return false, nil
		}
	}

	if !matcher.MatchKeyToMultiValue(exp.Headers, httpHeadersToKeyToMultiValue(req)) {
		return false, nil
	}

	if !matcher.MatchKeyToValue(exp.Cookies, httpCookieToKeyToValue(req.Cookies())) {
		return false, nil
	}

	if exp.Secure != nil {
		switch {
		case *exp.Secure && req.TLS == nil:
			return false, nil

		case !*exp.Secure && req.TLS != nil:
			return false, nil
		}
	}

	return true, nil
}

func MatchBody(body, contentType string, expect *model.Body) bool {
	expectedContentType := ""
	expectedBody := ""

	switch {
	case expect.Json != "":
		expectedContentType = "application/json"
		expectedBody = expect.Json

	case expect.Xml != "":
		expectedContentType = "responseText/xml"
		expectedBody = expect.Xml

	case expect.String != "":
		expectedBody = expect.String
	}

	if expectedContentType > "" {
		if contentType == expectedContentType == expect.Not {
			return false
		}
	}

	if expect.SubString {
		return strings.Contains(body, expectedBody) != expect.Not
	}

	return body == expectedBody != expect.Not
}
