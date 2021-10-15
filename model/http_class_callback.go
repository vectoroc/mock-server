package model

type HttpClassCallback struct {
	Delay         Delay  `json:"delay,omitempty"`
	CallbackClass string `json:"callbackClass,omitempty"`
}
