package server

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())
	now := time.Now()

	if cb, ok := s.apiRoutes[req.URL.Path]; ok && req.Method == "PUT" {
		defer func() {
			APIDurations.WithLabelValues(pathToMethod(req.URL.Path, s.apiPrefix)).Observe(time.Since(now).Seconds())
		}()

		cb(resp, req)
		return
	}
	defer func() {
		ProxyDurations.Observe(time.Since(now).Seconds())
	}()

	// match expectation
	expect, err := s.matchRequest(req)
	if err != nil {
		log.Err(err).Msg("internal error")
		internalError(resp, err)

		return
	}

	if expect == nil {
		log.Print("direct proxy request")
		if req.Method == "CONNECT" {
			ConnectMethod(resp, req)
			return
		}

		if req.URL.Host == "" {
			log.Error().Msg("empty host header")
			s.responseText(resp, http.StatusBadRequest, "URL expected")
			return
		}

		s.proxy.ServeHTTP(resp, req)
		return
	}

	switch {
	case expect.HttpResponse != nil:
		log.Print("replace response")
		ProcessHttpResponse(req.Context(), expect.HttpResponse, resp)
		return

	case expect.HttpError != nil:
		log.Print("replace response, error")
		ProcessHttpError(req.Context(), expect.HttpError, resp)
		return

	default:
		log.Error().Msg("not implemented")
		s.notImplemented(resp, req)
	}
}

func internalError(resp http.ResponseWriter, err error) {
	Codes.WithLabelValues(strconv.Itoa(http.StatusInternalServerError)).Inc()
	resp.WriteHeader(http.StatusInternalServerError)
	if _, err := resp.Write([]byte(err.Error())); err != nil {
		log.Err(err)
	}
}

func ConnectMethod(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())

	hij, ok := rw.(http.Hijacker)
	if !ok {
		log.Warn().Msg("http server does not support hijacker")
		return
	}

	clientConn, _, err := hij.Hijack()
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	proxyConn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	// The returned net.Conn may have read or write deadlines
	// already set, depending on the configuration of the
	// Server, to set or clear those deadlines as needed
	// we set timeout to 5 minutes
	deadline := time.Now().Add(time.Minute)

	err = clientConn.SetDeadline(deadline)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	err = proxyConn.SetDeadline(deadline)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	_, err = clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	go func() {
		_, err := io.Copy(clientConn, proxyConn)
		if err != nil {
			log.Printf("copy error: %s", err)
		}
		clientConn.Close()
		proxyConn.Close()
	}()

	if _, err := io.Copy(proxyConn, clientConn); err != nil {
		log.Printf("copy error: %s", err)
	}

	proxyConn.Close()
	clientConn.Close()
}

// pathToMethod returns API method name to be called
func pathToMethod(path string, prefix string) string {
	if strings.HasPrefix(path, prefix) {
		path = path[len(prefix):]
	}

	parts := strings.SplitN(strings.TrimLeft(path, "/"), "/", 1)
	switch len(parts) {
	case 0:
		return "unknown"
	default:
		return parts[0]
	}
}
