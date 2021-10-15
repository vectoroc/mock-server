package server

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"mock-server/matcher"
	"mock-server/model"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type Server struct {
	lock   *sync.RWMutex
	engine *matcher.Engine
	proxy  *httputil.ReverseProxy
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *Server {
	return &Server{
		lock: &sync.RWMutex{},
		proxy: &httputil.ReverseProxy{
			Director: func(request *http.Request) {},
			ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
				log.Err(err).CallerSkipFrame(4).Msg("proxy error")
			},
		},
		engine: matcher.NewEngine(),
		logger: logger,
	}
}

func notImplementedError(resp http.ResponseWriter) {
	log.Print("not implemented")
	resp.WriteHeader(http.StatusNotImplemented)
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

func ProcessHttpResponse(r *model.HttpResponse, w http.ResponseWriter) {
	ProccessDelay(r.Delay)

	for name, values := range r.Headers {
		for _, v := range values {
			w.Header().Add(name, v)
		}
	}

	for name, value := range r.Cookies {
		c := &http.Cookie{Name: name, Value: value}
		http.SetCookie(w, c)
	}

	switch {
	case r.Body.ContentType > "":
		w.Header().Add("Content-Type", r.Body.ContentType)

	case r.Body.Json > "":
		w.Header().Add("Content-Type", "application/json")

	case r.Body.Xml > "":
		w.Header().Add("Content-Type", "text/html")
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

	notImplementedError(w)
}

func ProcessHttpError(r *model.HttpError, w http.ResponseWriter) {
	ProccessDelay(r.Delay)

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

func ProccessDelay(d model.Delay) {
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

	time.Sleep(res)
}
