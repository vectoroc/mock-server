package server

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessHttpResponse(t *testing.T) {
	t.Run("default status code is 200", func(t *testing.T) {
		response := &model.HttpResponse{
			Body: model.BodyWithContentType{
				Type:   "text",
				String: `hello`,
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, 200, result.StatusCode)
		assert.Empty(t, result.Header)
		assert.Equal(t, response.Body.String, respRecorer.Body.String())
	})

	t.Run("non-2xx status code", func(t *testing.T) {
		response := &model.HttpResponse{
			StatusCode: 500,
			Body: model.BodyWithContentType{
				Type:   "text",
				String: `internal server error`,
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, 500, result.StatusCode)
		assert.Empty(t, result.Header)
		assert.Equal(t, response.Body.String, respRecorer.Body.String())
	})

	t.Run("redirect response", func(t *testing.T) {
		response := &model.HttpResponse{
			StatusCode: 302,
			Headers: model.KeyToMultiValue{
				Values: map[string][]string{
					"Location": {"https://www.test.com/redirect"},
				},
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, 302, result.StatusCode)
		assert.Equal(t, http.Header{
			"Location": {"https://www.test.com/redirect"},
		}, result.Header)
		assert.Equal(t, response.Body.String, respRecorer.Body.String())
	})

	t.Run("headers should be canonicalized", func(t *testing.T) {
		response := &model.HttpResponse{
			Headers: model.KeyToMultiValue{
				Values: map[string][]string{
					"x-myApp-Header": {"token12412421"},
				},
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, http.Header{
			"X-Myapp-Header": {"token12412421"},
		}, result.Header)
	})

	t.Run("json response", func(t *testing.T) {
		response := &model.HttpResponse{
			StatusCode: 200,
			Body: model.BodyWithContentType{
				Type: "JSON",
				Json: `{"x": {"y": [1,2,false]}}`,
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, 200, result.StatusCode)
		assert.Equal(t, http.Header{
			"Content-Type": {"application/json"},
		}, result.Header)
		assert.Equal(t, response.Body.Json, respRecorer.Body.String())
	})

	t.Run("cookies", func(t *testing.T) {
		response := &model.HttpResponse{
			Cookies: model.KeyToValue{
				Values: map[string]string{
					"sessionid": "some-secret",
				},
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, http.Header{
			"Set-Cookie": {"sessionid=some-secret"},
		}, result.Header)
	})

}
