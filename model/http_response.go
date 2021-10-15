package model

type HttpResponse struct {
	Delay             Delay               `json:"delay,omitempty"`
	Body              BodyWithContentType `json:"body,omitempty"`
	Cookies           KeyToValue          `json:"cookies,omitempty"`
	ConnectionOptions ConnectionOptions   `json:"connectionOptions,omitempty"`
	Headers           KeyToMultiValue     `json:"headers,omitempty"`
	StatusCode        int32               `json:"statusCode,omitempty"`
	ReasonPhrase      string              `json:"reasonPhrase,omitempty"`
}
