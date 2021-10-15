package model

type TimeToLive struct {
	TimeUnit   string `json:"timeUnit,omitempty"`
	TimeToLive int32  `json:"timeToLive,omitempty"`
	Unlimited  bool   `json:"unlimited,omitempty"`
}
