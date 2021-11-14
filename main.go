package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	stdlog "log"
	"mock-server/server"
	"net"
	"net/http"
	"os"
)

var (
	addr      = flag.String("api-addr", "127.0.0.1:8000", "")
	apiPrefix = flag.String("api-prefix", "/mockserver", "")
)

func main() {
	flag.Parse()

	logger := zerolog.New(os.Stderr)
	zerolog.DefaultContextLogger = &logger

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	s := server.New(logger, *apiPrefix)
	logger.Info().Str("addr", *addr).Msg("starting mock-server")

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Panic().Err(err).Msg("net.Listen has failed")
	}

	serv := http.Server{
		Handler: s.WrappedHandler(),
	}

	err = serv.Serve(l)
	if err != nil {
		logger.Panic().Err(err).Msg("unable to serve http")
	}
}
