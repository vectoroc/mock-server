package main

import (
	"context"
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
	addr = flag.String("api", "127.0.0.1:8000", "")
)

func main() {
	logger := zerolog.New(os.Stderr)

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	ctx := logger.WithContext(context.Background())

	s := server.New(logger)
	s.InitAPI()

	logger.Info().Str("add", *addr).Msg("starting mock-server")

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	serv := http.Server{
		Handler: hlog.RequestHandler("proxy")(s),
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	err = serv.Serve(l)
	if err != nil {
		panic(err)
	}
}
