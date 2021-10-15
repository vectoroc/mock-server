package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mock-server/model"
	"net/http"
)

func (s *Server) Expectation(req *model.Expectations) {
	s.engine.AddExpectations(req.ToArray())
}

func (s *Server) Clear(exp *model.HttpRequest) {
	s.engine.ClearBy(exp)
}

func (s *Server) Reset() {
	s.engine.Reset()
}

func (s *Server) validateAPIRequest(context *gin.Context, request interface{}) bool {
	if context.ContentType() != "application/json" {
		context.String(http.StatusBadRequest, "incorrect request format")

		return false
	}
	defer context.Request.Body.Close()

	err := json.NewDecoder(context.Request.Body).Decode(request)
	if err != nil {
		s.logger.Printf("failed to unmarshal body %s", err)
		context.String(http.StatusNotAcceptable, "invalid expectation")

		return false
	}

	return true
}
