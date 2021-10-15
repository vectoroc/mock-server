package model

// Delay - response delay
type Delay struct {
	// - DAYS
	// - HOURS
	// - MINUTES
	// - SECONDS
	// - MILLISECONDS
	// - MICROSECONDS
	// - NANOSECONDS
	TimeUnit string `json:"timeUnit,omitempty"`
	Value    int32  `json:"value,omitempty"`
}
