package model

type HttpForward struct {
	Delay  Delay  `json:"delay,omitempty"`
	Host   string `json:"host,omitempty"`
	Port   int32  `json:"port,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}
