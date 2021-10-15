package model

type HttpRequest struct {
	Body                  *Body           `json:"body,omitempty"`
	Headers               KeyToMultiValue `json:"headers,omitempty"`
	Cookies               KeyToValue      `json:"cookies,omitempty"`
	QueryStringParameters KeyToMultiValue `json:"queryStringParameters,omitempty"`
	Path                  string          `json:"path,omitempty"`
	Method                string          `json:"method,omitempty"`
	Secure                *bool           `json:"secure,omitempty"`
	KeepAlive             *bool           `json:"keepAlive,omitempty"`
	// SocketAddress         *SocketAddress  `json:"socketAddress,omitempty"`
}

func NewHttpRequest() HttpRequest {
	return HttpRequest{
		// preallocate custom container types otherwise UnmarshalJSON will panic
		Headers:               KeyToMultiValue{},
		Cookies:               KeyToValue{},
		QueryStringParameters: KeyToMultiValue{},
	}
}