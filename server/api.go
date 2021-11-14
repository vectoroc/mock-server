package server

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"mock-server/model"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) InitAPI() {
	s.apiRoutes = map[string]http.HandlerFunc{
		s.apiPrefix + "/expectation":    s.Expectation,
		s.apiPrefix + "/clear":          s.Clear,
		s.apiPrefix + "/reset":          s.Reset,
		s.apiPrefix + "/retrieve":       s.Retrieve,
		s.apiPrefix + "/verify":         s.notImplemented,
		s.apiPrefix + "/verifySequence": s.notImplemented,
		s.apiPrefix + "/responseStatus": s.notImplemented,
		s.apiPrefix + "/bind":           s.notImplemented,
		s.apiPrefix + "/stop":           s.notImplemented,
	}
}

func (s *Server) Expectation(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())

	params := &model.Expectations{}
	if !s.validateAPIRequest(rw, req, params) {
		return
	}
	list := params.ToArray()
	s.engine.AddExpectations(list)
	log.Info().Int("added", len(list)).Send()

	ExpectationsAdd.Add(float64(len(list)))
	s.responseJSON(rw, http.StatusCreated, list)
}

func (s *Server) Clear(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())

	exp := &model.HttpRequest{}
	if !s.validateAPIRequest(rw, req, exp) {
		return
	}
	ids := s.engine.ClearBy(exp)
	log.Info().Int("expectations", len(ids)).Send()

	ExpectationsClear.Add(float64(len(ids)))
	s.responseJSON(rw, http.StatusOK, ids)
}

func (s *Server) Reset(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())

	s.engine.Reset()
	log.Info().Send()
	s.responseStatus(rw, http.StatusOK)
}

func (s *Server) Retrieve(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())
	log.Info().Send()

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

func (s *Server) reqContentType(req *http.Request) string {
	header := req.Header.Get("Content-Type")
	return strings.TrimSpace(strings.Split(header, ";")[0])
}

func (s *Server) validateAPIRequest(rw http.ResponseWriter, req *http.Request, params interface{}) bool {
	log := zerolog.Ctx(req.Context())

	ct := s.reqContentType(req)
	if ct != "application/json" {
		log.Warn().Str("content-type", ct).Msg("incorrect request format")
		s.responseText(rw, http.StatusBadRequest, "incorrect request format")
		return false
	}

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal body")
		s.responseText(rw, http.StatusNotAcceptable, "invalid expectation")
		return false
	}

	return true
}

func (s *Server) responseJSON(rw http.ResponseWriter, status int, data interface{}) {
	s.responseStatus(rw, status)
	enc := json.NewEncoder(rw)
	enc.SetIndent("", " ")
	err := enc.Encode(data)
	if err != nil {
		s.logger.Err(err).Msg("failed to marshal json response")
	}
}

func (s *Server) responseText(rw http.ResponseWriter, status int, resp string) {
	s.responseStatus(rw, status)
	if _, err := rw.Write([]byte(resp)); err != nil {
		s.logger.Err(err).Msg("failed to send json response")
	}
}

func (s *Server) responseStatus(rw http.ResponseWriter, status int) {
	Codes.WithLabelValues(strconv.Itoa(status)).Inc()
	rw.WriteHeader(status)
}

func (s *Server) notImplemented(rw http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())
	log.Warn().Msg("request handler is not implemented")
	s.responseStatus(rw, http.StatusNotImplemented)
}
