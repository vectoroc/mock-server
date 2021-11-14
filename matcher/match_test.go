package matcher

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"testing"
)

func TestMatchKeyToMultiValue(t *testing.T) {
	t.Run("it should match only headers present in expectation", func(t *testing.T) {
		exp := model.KeyToMultiValue{
			Values: map[string][]string{
				"Host": {"ya.ru"},
			},
		}

		headers := model.KeyToMultiValue{
			Values: map[string][]string{
				"Accept":           {"*/*"},
				"Host":             {"ya.ru"},
				"Proxy-Connection": {"Keep-Alive"},
				"User-Agent":       {"curl/7.64.1"},
			},
		}

		assert.True(t, MatchKeyToMultiValue(exp, headers))

		t.Run("extra headers in expectation should fail", func(t *testing.T) {
			exp.Values["User-Agent"] = []string{"chrome"}
			assert.False(t, MatchKeyToMultiValue(exp, headers))
		})
	})

	t.Run("values should fully match", func(t *testing.T) {
		exp := model.KeyToMultiValue{
			Values: map[string][]string{
				"content-type": {"text/plain"},
			},
		}

		headers := model.KeyToMultiValue{
			Values: map[string][]string{
				"content-type": {"text/plain", "text/csv"},
			},
		}

		assert.False(t, MatchKeyToMultiValue(exp, headers))
	})
}

func TestMatchKeyToValue(t *testing.T) {
	exp := model.KeyToValue{
		Values: map[string]string{
			"domain":     "ya.ru",
			"session_id": "123",
		},
	}

	cookies := model.KeyToValue{
		Values: map[string]string{
			"domain":     "ya.ru",
			"token":      "2412515151251521",
			"session_id": "123",
		},
	}

	assert.True(t, MatchKeyToValue(exp, cookies))

	exp.Values["token"] = ""
	assert.False(t, MatchKeyToValue(exp, cookies), "all expected fields should match")
}

func TestMatchRequestByRequest(t *testing.T) {
	expectation := model.HttpRequest{
		Method: "GET",
		Path:   "/test",
	}
	pattern1 := model.HttpRequest{
		Method: "GET",
		Path:   "/test",
	}
	pattern2 := model.HttpRequest{
		Method: "POST",
		Path:   "/test",
	}
	pattern3 := model.HttpRequest{
		Method: "GET",
		Path:   "/test/1",
	}
	pattern4 := model.HttpRequest{
		Method: "POST",
		Path:   "/test",
	}
	assert.True(t, MatchRequestByRequest(&expectation, &pattern1))
	assert.False(t, MatchRequestByRequest(&expectation, &pattern2))
	assert.False(t, MatchRequestByRequest(&expectation, &pattern3))
	assert.False(t, MatchRequestByRequest(&expectation, &pattern4))
}
