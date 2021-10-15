package model

type Expectation struct {
	Id                           string                        `json:"id,omitempty"`
	Priority                     int                           `json:"priority,omitempty"`
	HttpRequest                  HttpRequest                   `json:"httpRequest,omitempty"`
	HttpResponse                 *HttpResponse                 `json:"httpResponse,omitempty"`
	HttpResponseTemplate         *HttpTemplate                 `json:"httpResponseTemplate,omitempty"`
	HttpResponseClassCallback    *HttpClassCallback            `json:"httpResponseClassCallback,omitempty"`
	HttpResponseObjectCallback   *HttpObjectCallback           `json:"httpResponseObjectCallback,omitempty"`
	HttpForward                  *HttpForward                  `json:"httpForward,omitempty"`
	HttpForwardTemplate          *HttpTemplate                 `json:"httpForwardTemplate,omitempty"`
	HttpForwardClassCallback     *HttpClassCallback            `json:"httpForwardClassCallback,omitempty"`
	HttpForwardObjectCallback    *HttpObjectCallback           `json:"httpForwardObjectCallback,omitempty"`
	HttpOverrideForwardedRequest *HttpOverrideForwardedRequest `json:"httpOverrideForwardedRequest,omitempty"`
	HttpError                    *HttpError                    `json:"httpError,omitempty"`  // error behaviour
	Times                        *Times                        `json:"times,omitempty"`      // number of responses
	TimeToLive                   *TimeToLive                   `json:"timeToLive,omitempty"` // time expectation is valid for
}
