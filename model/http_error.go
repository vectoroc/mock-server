package model

type HttpError struct {
	Delay          Delay  `json:"delay,omitempty"`
	DropConnection bool   `json:"dropConnection,omitempty"`
	ResponseBytes  string `json:"responseBytes,omitempty"`
}
