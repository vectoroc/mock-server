package server

import (
	"github.com/gin-gonic/gin"
	"mock-server/model"
	"net/http"
)

func (s *Server) InitRoutes(r *gin.Engine) error {

	r.PUT("/expectation", func(context *gin.Context) {
		req := &model.Expectations{}
		if !s.validateAPIRequest(context, req) {
			return
		}

		s.Expectation(req)
		context.JSON(http.StatusCreated, req)
	})

	r.PUT("/clear", func(context *gin.Context) {
		req := model.NewHttpRequest()
		if !s.validateAPIRequest(context, &req) {
			return
		}

		s.Clear(&req)
		context.Status(http.StatusOK)
	})

	r.PUT("/reset", func(context *gin.Context) {
		s.Reset()
		context.Status(http.StatusOK)
	})

	r.PUT("/retrieve", func(context *gin.Context) {
		format := context.Query("format")
		sType := context.Query("type")

		if format == "" {
			format = "json"
		}

		if format != "json" {
			notImplementedError(context.Writer)
			return
		}

		//-logs
		//-requests
		//-request_responses
		//-recorded_expectations
		//-active_expectations
		switch sType {
		case "active_expectations":
			context.JSON(http.StatusOK, s.engine.Expectations())
			return

		default:
			notImplementedError(context.Writer)
			return
		}
	})

	r.PUT("/verify", func(context *gin.Context) {
		context.Status(http.StatusNotImplemented)
	})

	r.PUT("/verifySequence", func(context *gin.Context) {
		context.Status(http.StatusNotImplemented)
	})

	r.PUT("/status", func(context *gin.Context) {
		context.Status(http.StatusNotImplemented)
	})

	r.PUT("/bind", func(context *gin.Context) {
		context.Status(http.StatusNotImplemented)
	})

	r.PUT("/stop", func(context *gin.Context) {
		context.Status(http.StatusNotImplemented)
	})
	return nil
}
