package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	stdlog "log"
	"mock-server/server"
	"net"
	"net/http"
	"os"
)

var (
	proxyAddr = flag.String("api", "127.0.0.1:8000", "")
	apiAddr   = flag.String("proxy", "127.0.0.1:8001", "")
)

func main() {
	logger := zerolog.New(os.Stderr)

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	ctx := logger.WithContext(context.Background())

	r := gin.Default()

	s := server.New(logger)

	if err := s.InitRoutes(r); err != nil {
		panic(err)
	}

	go proxyListener(ctx, s, *proxyAddr)
	initAPIListener(ctx, r, *apiAddr)
}

func initAPIListener(ctx context.Context, r *gin.Engine, addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	if err := r.RunListener(l); err != nil {
		panic(err)
	}
}

func proxyListener(ctx context.Context, s http.Handler, addr string) {
	l, err := net.Listen("tcp", addr)
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
