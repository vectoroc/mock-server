package server

import (
	"encoding/json"
	"mock-server/model"
	"net/http"
)

func (s *Server) InitAPI() {
	s.apiRoutes = map[string]http.HandlerFunc{
		"/expectation":    s.Expectation,
		"/clear":          s.Clear,
		"/reset":          s.Reset,
		"/retrieve":       s.Retrieve,
		"/verify":         s.notImplemented,
		"/verifySequence": s.notImplemented,
		"/responseStatus": s.notImplemented,
		"/bind":           s.notImplemented,
		"/stop":           s.notImplemented,
	}
}

func (s *Server) Expectation(rw http.ResponseWriter, req *http.Request) {
	params := &model.Expectations{}
	if !s.validateAPIRequest(rw, req, &params) {
		return
	}
	list := params.ToArray()
	s.engine.AddExpectations(list)
	s.responseJSON(rw, http.StatusCreated, list)
}

func (s *Server) Clear(rw http.ResponseWriter, req *http.Request) {
	exp := &model.HttpRequest{}
	if !s.validateAPIRequest(rw, req, &exp) {
		return
	}
	s.engine.ClearBy(exp)
	s.responseStatus(rw, http.StatusOK)
}

func (s *Server) Reset(rw http.ResponseWriter, req *http.Request) {
	s.engine.Reset()
	s.responseStatus(rw, http.StatusOK)
}

func (s *Server) Retrieve(rw http.ResponseWriter, req *http.Request) {
	format := req.URL.Query().Get("format")
	sType := req.URL.Query().Get("type")

	if format == "" {
		format = "json"
	}

	if format != "json" {
		s.notImplemented(rw, req)
		return
	}

	//-logs
	//-requests
	//-request_responses
	//-recorded_expectations
	//-active_expectations
	switch sType {
	case "active_expectations":
		s.responseJSON(rw, http.StatusOK, s.engine.Expectations())
		return

	default:
		s.notImplemented(rw, req)
		return
	}
}

func (s *Server) validateAPIRequest(rw http.ResponseWriter, req *http.Request, params interface{}) bool {
	ct := req.Header.Get("Content-Type")
	if ct != "application/json" {
		s.logger.Warn().Str("content-type", ct).Msg("incorrect request format")
		s.responseText(rw, http.StatusBadRequest, "incorrect request format")
		return false
	}

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		s.logger.Err(err).Msg("failed to unmarshal body")
		s.responseText(rw, http.StatusNotAcceptable, "invalid expectation")
		return false
	}

	return true
}

func (s *Server) responseJSON(rw http.ResponseWriter, status int, data interface{}) {
	rw.WriteHeader(status)
	resp, err := json.Marshal(data)
	if err != nil {
		s.logger.Err(err).Msg("failed to marshal json response")
	}
	if _, err := rw.Write(resp); err != nil {
		s.logger.Err(err).Msg("failed to send json response")
	}
}

func (s *Server) responseText(rw http.ResponseWriter, status int, resp string) {
	rw.WriteHeader(status)
	if _, err := rw.Write([]byte(resp)); err != nil {
		s.logger.Err(err).Msg("failed to send json response")
	}
}

func (s *Server) responseStatus(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
}

func (s *Server) notImplemented(rw http.ResponseWriter, req *http.Request) {
	s.responseStatus(rw, http.StatusNotImplemented)
}
