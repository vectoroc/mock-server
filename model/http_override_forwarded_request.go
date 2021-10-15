package model

type HttpOverrideForwardedRequest struct {
	Delay        Delay        `json:"delay,omitempty"`
	HttpRequest  HttpRequest  `json:"httpRequest,omitempty"`
	HttpResponse HttpResponse `json:"httpResponse,omitempty"`
}
