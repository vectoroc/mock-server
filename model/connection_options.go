package model

// ConnectionOptions - connection options
type ConnectionOptions struct {
	CloseSocket                 bool  `json:"closeSocket,omitempty"`
	CloseSocketDelay            Delay `json:"closeSocketDelay,omitempty"`
	ContentLengthHeaderOverride int32 `json:"contentLengthHeaderOverride,omitempty"`
	SuppressContentLengthHeader bool  `json:"suppressContentLengthHeader,omitempty"`
	SuppressConnectionHeader    bool  `json:"suppressConnectionHeader,omitempty"`
	KeepAliveOverride           bool  `json:"keepAliveOverride,omitempty"`
}
