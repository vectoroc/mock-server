package matcher

import (
	"mock-server/model"
)

func MatchRequestByRequest(exp *model.HttpRequest, req *model.HttpRequest) bool {
	if exp.Method != "" && req.Method != exp.Method {
		return false
	}

	if exp.Path > "" && exp.Path != req.Path {
		return false
	}

	if exp.Secure != nil && (req.Secure == nil || *exp.Secure != *req.Secure) {
		return false
	}

	if exp.KeepAlive != nil && (req.KeepAlive == nil || *exp.KeepAlive != *req.KeepAlive) {
		return false
	}

	if !MatchBodyByBody(exp.Body, req.Body) {
		return false
	}

	if !MatchKeyToMultiValue(exp.Headers, req.Headers) {
		return false
	}

	if !MatchKeyToValue(exp.Cookies, req.Cookies) {
		return false
	}

	if !MatchKeyToMultiValue(exp.QueryStringParameters, req.QueryStringParameters) {
		return false
	}

	return true
}

func MatchBodyByBody(exp *model.Body, body *model.Body) bool {
	if exp == nil {
		return true
	}

	// bool fields
	if exp.Not && exp.Not != body.Not {
		return false
	}
	if exp.SubString && exp.SubString != body.SubString {
		return false
	}

	// string fields
	if exp.Type > "" && exp.Type != body.Type {
		return false
	}
	if exp.Base64Bytes > "" && exp.Base64Bytes != body.Base64Bytes {
		return false
	}
	if exp.ContentType > "" && exp.ContentType != body.ContentType {
		return false
	}
	if exp.Json > "" && exp.Json != body.Json {
		return false
	}
	if exp.MatchType > "" && exp.MatchType != body.MatchType {
		return false
	}
	if exp.JsonSchema > "" && exp.JsonSchema != body.JsonSchema {
		return false
	}
	if exp.JsonPath > "" && exp.JsonPath != body.JsonPath {
		return false
	}
	if exp.Regex > "" && exp.Regex != body.Regex {
		return false
	}
	if exp.String > "" && exp.String != body.String {
		return false
	}
	if exp.Xml > "" && exp.Xml != body.Xml {
		return false
	}
	if exp.XmlSchema > "" && exp.XmlSchema != body.XmlSchema {
		return false
	}
	if exp.Xpath > "" && exp.Xpath != body.Xpath {
		return false
	}

	return true
}

func MatchKeyToMultiValue(expect model.KeyToMultiValue, actual model.KeyToMultiValue) bool {
	if len(expect) == 0 {
		return true
	}

	if len(expect) > len(actual) {
		return false
	}

	for name, values := range expect {
		if !Equal(values, actual[name]) {
			return false
		}
	}

	return true
}

func MatchKeyToValue(expect model.KeyToValue, actual model.KeyToValue) bool {
	if len(expect) == 0 {
		return true
	}

	if len(expect) > len(actual) {
		return false
	}

	for name, v := range expect {
		c, ok := actual[name]
		if !ok || c != v {
			return false
		}
	}

	return true
}
