package model

type HttpObjectCallback struct {
	Delay            Delay  `json:"delay,omitempty"`
	ClientId         string `json:"clientId,omitempty"`
	ResponseCallback bool   `json:"responseCallback,omitempty"`
}
