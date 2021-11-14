package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ExpectationsAdd = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mockserver",
		Subsystem: "expectations",
		Name:      "add",
		Help:      "Added expectations counter",
	})

	ExpectationsClear = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mockserver",
		Subsystem: "expectations",
		Name:      "clear",
		Help:      "Cleared expectations counter",
	})

	Codes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "mockserver",
		Subsystem: "api",
		Name:      "codes",
		Help:      "Api status codes",
	}, []string{"code"})

	APIDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "mockserver",
		Subsystem: "api",
		Name:      "duration_seconds",
		Help:      "API calls durations by method",
	}, []string{"method"})

	ProxyDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "mockserver",
		Subsystem: "proxy",
		Name:      "duration_seconds",
		Help:      "API calls durations by method",
	})
)

func init() {
	prometheus.MustRegister(ExpectationsAdd)
	prometheus.MustRegister(ExpectationsClear)
	prometheus.MustRegister(APIDurations)
}
