package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"mock-server/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestProcessHttpResponse(t *testing.T) {
	ctx := context.Background()
	t.Run("default status code is 200", func(t *testing.T) {
		response := &model.HttpResponse{
			Body: model.BodyWithContentType{
				Type:   "text",
				String: `hello`,
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(ctx, response, respRecorer)

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
		ProcessHttpResponse(ctx, response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, 500, result.StatusCode)
		assert.Empty(t, result.Header)
		assert.Equal(t, response.Body.String, respRecorer.Body.String())
	})

	t.Run("redirect response", func(t *testing.T) {
		response := &model.HttpResponse{
			StatusCode: 302,
			Headers: model.KeyToMultiValue{
				"Location": {"https://www.test.com/redirect"},
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(ctx, response, respRecorer)

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
				"x-myApp-Header": {"token12412421"},
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(ctx, response, respRecorer)

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
		ProcessHttpResponse(ctx, response, respRecorer)

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
				"sessionid": "some-secret",
			},
		}

		respRecorer := httptest.NewRecorder()
		ProcessHttpResponse(ctx, response, respRecorer)

		result := respRecorer.Result()

		assert.Equal(t, http.Header{
			"Set-Cookie": {"sessionid=some-secret"},
		}, result.Header)
	})

}

func TestServer(t *testing.T) {
	contentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if val := r.URL.Query().Get("delay"); val != "" {
			d, _ := strconv.Atoi(val)
			if d > 0 {
				select {
				case <-time.NewTimer(time.Millisecond * time.Duration(d)).C:
				case <-r.Context().Done():
					return
				}
			}
		}
		_, err := w.Write([]byte("HELLO " + r.URL.String()))
		if err != nil {
			log.Err(err)
		}
	}))
	t.Cleanup(contentServer.Close)

	s := New(zerolog.Nop(), "/api")
	mockServer := httptest.NewServer(s.WrappedHandler())
	t.Cleanup(mockServer.Close)

	proxyUrl, err := url.Parse(mockServer.URL)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	t.Parallel()

	initWG, clearWG := sync.WaitGroup{}, sync.WaitGroup{}

	t.Run("init 200 ok expecations", func(t *testing.T) {
		initWG.Add(1)
		defer initWG.Done()
		for i := 0; i < 100; i++ {
			initSimpleExpectation(t, mockServer.URL+"/api", i)
		}
	})

	t.Run("init redirect expectations", func(t *testing.T) {
		initWG.Add(1)
		defer initWG.Done()
		for i := 0; i < 100; i++ {
			initRedirectExpectation(t, mockServer.URL+"/api", i)
		}
	})

	t.Run("init error expectations", func(t *testing.T) {
		initWG.Add(1)
		defer initWG.Done()
		for i := 400; i < 500; i++ {
			initErrorExpectation(t, mockServer.URL+"/api", i)
		}
	})

	// run requests
	t.Run("simple proxy request", func(t *testing.T) {
		wg := sync.WaitGroup{}
		ch := make(chan string, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for u := range ch {
					resp, err := client.Get(contentServer.URL + u)
					require.NoError(t, err)
					assert.Equal(t, 200, resp.StatusCode)

					body, err := io.ReadAll(resp.Body)
					require.NoError(t, err)
					require.NoError(t, resp.Body.Close())

					assert.Equal(t, "HELLO "+u, string(body))
				}
			}()
		}

		for i := 0; i < 10000; i++ {
			delay := i%50 + 1
			if i%100 == 0 {
				delay = 10000
			}
			ch <- "/some/simple/request?foo=bar&delay=" + strconv.Itoa(delay)
		}

		close(ch)
		wg.Wait()
	})

	t.Run("mocked 200-ok requests", func(t *testing.T) {
		initWG.Wait() // wait until all expectations set

		for i := 0; i < 100; i++ {
			resp, err := client.Get(fmt.Sprintf("%s/test/%d", contentServer.URL, i))
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "simple-expectation", resp.Header.Get("x-mock-server"))
			assert.Equal(t, "test mock server", string(body))
		}
	})

	t.Run("mocked redirect requests", func(t *testing.T) {
		clearWG.Add(1)
		defer clearWG.Done()
		initWG.Wait() // wait until all expectations set

		for i := 0; i < 100; i++ {
			resp, err := client.Get(fmt.Sprintf("%s/test-redirect/%d", contentServer.URL, i))
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "test mock server", string(body))
		}
	})

	t.Run("mocked error requests", func(t *testing.T) {
		clearWG.Add(1)
		defer clearWG.Done()
		initWG.Wait() // wait until all expectations set

		for i := 400; i < 500; i++ {
			resp, err := client.Get(fmt.Sprintf("%s/test-error/%d", contentServer.URL, i))
			require.NoError(t, err)
			assert.Equal(t, i, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "test error", string(body))
		}
	})

	t.Run("clear expectations", func(t *testing.T) {
		clearWG.Wait()
		for i := 0; i < 100; i++ {
			exp := fmt.Sprintf(`{
    "method" : "GET",
    "path" : "/test/%d"
}`, i)

			ids := clearExpectation(t, mockServer.URL+"/api", strings.NewReader(exp))
			assert.Len(t, ids, 1)
		}
	})

	t.Run("delay check", func(t *testing.T) {
		newDelayExp := fmt.Sprintf(`{
  "httpRequest": {
    "method": "GET",
    "path": "/delay"
  },
  "httpResponse": {
    "statusCode": 200,
    "body": "test delay",
    "delay": {
      "timeUnit": "SECONDS",
      "value": 10
    }
  }
}`)

		addExpectation(t, mockServer.URL+"/api", strings.NewReader(newDelayExp))

		ts := time.Now()
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", contentServer.URL+"/delay", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		assert.Error(t, err)
		assert.Empty(t, resp)
		assert.Less(t, time.Since(ts), time.Second*2)
	})
}

func addExpectation(t *testing.T, mockServerURL string, expectation io.Reader) model.Expectations {
	req, err := http.NewRequest("PUT", mockServerURL+"/expectation", expectation)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, 201, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	if resp.StatusCode != 201 {
		t.Error(string(respBody))
	}

	var expectations model.Expectations
	err = json.Unmarshal(respBody, &expectations)
	require.NoError(t, err)

	return expectations
}

func clearExpectation(t *testing.T, mockServerURL string, expectation io.Reader) (removedIds []string) {
	req, err := http.NewRequest("PUT", mockServerURL+"/clear", expectation)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, 200, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	if resp.StatusCode != 200 {
		t.Error(string(respBody))
	}

	err = json.Unmarshal(respBody, &removedIds)
	require.NoError(t, err)
	return
}

func initSimpleExpectation(t *testing.T, mockServerURL string, i int) {
	newExp := fmt.Sprintf(`{
  "httpRequest": {
    "method": "GET",
    "path": "/test/%d"
  },
  "httpResponse": {
    "statusCode": 200,
    "body": "test mock server",
    "headers": {
      "x-mock-server": "simple-expectation"
    }
  }
}`, i)

	expectations := addExpectation(t, mockServerURL, strings.NewReader(newExp))
	assert.Len(t, expectations.ToArray(), 1)
}

func initRedirectExpectation(t *testing.T, mockServerURL string, i int) {
	newExp := fmt.Sprintf(`{
  "httpRequest" : {
    "method" : "GET",
    "path" : "/test-redirect/%d"
  },
  "httpResponse" : {
    "statusCode" : 301,
    "headers": {
      "location": "/test/%d"
    }
  }
}`, i, i)

	expectations := addExpectation(t, mockServerURL, strings.NewReader(newExp))
	assert.Len(t, expectations.ToArray(), 1)
}

func initErrorExpectation(t *testing.T, mockServerURL string, code int) {
	newExp := fmt.Sprintf(`{
  "httpRequest" : {
    "method" : "GET",
    "path" : "/test-error/%d"
  },
  "httpResponse" : {
    "statusCode" : %d,
    "body": "test error"
  }
}`, code, code)

	expectations := addExpectation(t, mockServerURL, strings.NewReader(newExp))
	assert.Len(t, expectations.ToArray(), 1)
}
