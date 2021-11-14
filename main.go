package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	stdlog "log"
	"mock-server/server"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

var (
	addr      = flag.String("api-addr", "127.0.0.1:8000", "")
	apiPrefix = flag.String("api-prefix", "/mockserver", "")

	metrics              = flag.String("metrics", "0.0.0.0:8001", "")
	mutexProfileFraction = flag.Int("mutex-profile-fraction", 0, "")
	blockProfileRatio    = flag.Int("block-profile-ratio", 0, "")
)

func main() {
	flag.Parse()

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger.With().Caller().Logger())

	s := server.New(logger, *apiPrefix)
	logger.Info().Str("addr", *addr).Msg("starting mock-server")

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Panic().Err(err).Msg("net.Listen has failed")
	}

	serv := http.Server{
		Handler: s.WrappedHandler(),
		//IdleTimeout: time.Minute,
	}

	if *blockProfileRatio > 0 {
		logger.Info().Int("ratio", *blockProfileRatio).Msg("enable blockProfileRatio")
		runtime.SetBlockProfileRate(*blockProfileRatio)
	}

	if *mutexProfileFraction > 0 {
		logger.Info().Int("ratio", *mutexProfileFraction).Msg("enable mutexProfileFraction")
		runtime.SetMutexProfileFraction(*mutexProfileFraction)
	}

	if *metrics > "" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			logger.Info().Str("addr", *metrics).Msg("listen metrics handler")
			if err := http.ListenAndServe(*metrics, nil); err != nil {
				logger.Panic().Err(err).Send()
			}
		}()
	}

	err = serv.Serve(l)
	if err != nil {
		logger.Panic().Err(err).Msg("unable to serve http")
	}
}
