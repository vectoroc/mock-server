package model

type Times struct {
	RemainingTimes int32 `json:"remainingTimes,omitempty"`
	Unlimited      bool  `json:"unlimited,omitempty"`
}
