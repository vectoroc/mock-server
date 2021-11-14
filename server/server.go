package server

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"mock-server/matcher"
	"mock-server/model"
	"net/http"
	"net/http/httputil"
	"time"
)

type Server struct {
	engine matcher.Engine
	proxy  httputil.ReverseProxy
	logger zerolog.Logger

	apiPrefix string
	apiRoutes map[string]http.HandlerFunc

	m []middleware
}

func New(logger zerolog.Logger, apiPrefix string) *Server {
	s := &Server{
		proxy: httputil.ReverseProxy{
			Director: func(request *http.Request) {},
			ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
				log.Ctx(request.Context()).Err(err).Msg("proxy error")
			},
		},
		logger:    logger,
		apiPrefix: apiPrefix,
	}
	s.InitAPI()
	s.Middleware(hlog.NewHandler(logger))
	s.Middleware(hlog.RequestIDHandler("id", ""))
	s.Middleware(hlog.AccessHandler(accessHandler))
	return s
}

func accessHandler(r *http.Request, status, size int, duration time.Duration) {
	hlog.FromRequest(r).Info().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Int("status", status).
		Int("size", size).
		Dur("duration", duration).
		Str("user-agent", r.UserAgent()).
		Str("remote-addr", r.RemoteAddr).
		Send()
}

func (s *Server) matchRequest(req *http.Request) (*model.Expectation, error) {

	for _, exp := range s.engine.Expectations() {
		ok, err := matchHttpRequest(&exp.HttpRequest, req)
		log.Ctx(req.Context()).Printf("expectation %s match %v err=%v", exp.Id, ok, err)
		if err != nil {
			return nil, err
		}
		if ok {
			return exp, nil
		}
	}

	return nil, nil
}

func ProcessHttpResponse(ctx context.Context, r *model.HttpResponse, w http.ResponseWriter) {
	ProccessDelay(ctx, r.Delay)
	if ctx.Err() != nil {
		return
	}

	for name, values := range r.Headers.Values {
		for _, v := range values {
			w.Header().Add(name, v)
		}
	}

	for name, value := range r.Cookies.Values {
		c := &http.Cookie{Name: name, Value: value}
		http.SetCookie(w, c)
	}

	switch {
	case r.Body.ContentType > "":
		w.Header().Add("Content-Type", r.Body.ContentType)

	case r.Body.Json > "":
		w.Header().Add("Content-Type", "application/json")

	case r.Body.Xml > "":
		w.Header().Add("Content-Type", "responseText/html")
	}

	status := http.StatusOK
	if r.StatusCode > 0 {
		status = int(r.StatusCode)
	}
	w.WriteHeader(status)

	switch {
	case r.Body.String > "":
		_, err := w.Write([]byte(r.Body.String))
		if err != nil {
			internalError(w, err)
		}
		return

	case r.Body.Json > "":
		_, err := w.Write([]byte(r.Body.Json))
		if err != nil {
			internalError(w, err)
		}
		return

	case r.Body.Xml > "":
		_, err := w.Write([]byte(r.Body.Xml))
		if err != nil {
			internalError(w, err)
		}
		return
	}
}

func ProcessHttpError(ctx context.Context, r *model.HttpError, w http.ResponseWriter) {
	ProccessDelay(ctx, r.Delay)
	if ctx.Err() != nil {
		return
	}

	if r.ResponseBytes > "" {
		_, err := w.Write([]byte(r.ResponseBytes))
		if err != nil {
			log.Err(err)
		}
	}
	if r.DropConnection {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			log.Print("not supported")
			return
		}

		conn, _, err := hijacker.Hijack()
		if err != nil {
			log.Err(err)
			return
		}

		if err := conn.Close(); err != nil {
			log.Err(err)
		}
	}
}

func ProccessDelay(ctx context.Context, d model.Delay) {
	if d.Value == 0 {
		return
	}

	res := time.Duration(d.Value)

	switch d.TimeUnit {
	case "DAYS":
		res *= time.Hour * 24
	case "HOURS":
		res *= time.Hour
	case "MINUTES":
		res *= time.Minute
	case "SECONDS":
		res *= time.Second
	case "MILLISECONDS":
		res *= time.Millisecond
	case "MICROSECONDS":
		res *= time.Microsecond
	case "NANOSECONDS":
		res *= time.Nanosecond
	}

	t := time.NewTimer(res)
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
