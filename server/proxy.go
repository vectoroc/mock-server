package server

import (
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"net/http"
	"time"
)

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log := hlog.FromRequest(req)

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
		} else {
			s.proxy.ServeHTTP(resp, req)
		}

		return
	}

	switch {
	case expect.HttpResponse != nil:
		log.Print("replace response")
		ProcessHttpResponse(expect.HttpResponse, resp)
		return

	case expect.HttpError != nil:
		log.Print("replace response, error")
		ProcessHttpError(expect.HttpError, resp)
		return

	default:
		log.Error().Msg("not implemented")
		notImplementedError(resp)
	}
}

func internalError(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusInternalServerError)
	if _, err := resp.Write([]byte(err.Error())); err != nil {
		log.Err(err)
	}
}

func ConnectMethod(rw http.ResponseWriter, req *http.Request) {
	log := hlog.FromRequest(req)

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
