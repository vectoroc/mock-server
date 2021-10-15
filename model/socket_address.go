package model

type SocketAddress struct {
	Host   string `json:"host,omitempty"`
	Port   int32  `json:"port,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}
