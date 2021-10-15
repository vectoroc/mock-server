package server

import "net/http"

type middleware func(http.Handler) http.Handler

func (s *Server) Middleware(m middleware) {
	s.m = append(s.m, m)
}

func (s *Server) WrappedHandler() http.Handler {
	var h http.Handler = s
	for i := range s.m {
		h = s.m[len(s.m)-1-i](h)
	}
	return h
}
