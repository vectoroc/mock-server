package matcher

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"testing"
)

func TestMatchKeyToMultiValue(t *testing.T) {
	exp := model.KeyToMultiValue{
		"Host": {"ya.ru"},
	}

	headers := model.KeyToMultiValue{
		"Accept":           {"*/*"},
		"Host":             {"ya.ru"},
		"Proxy-Connection": {"Keep-Alive"},
		"User-Agent":       {"curl/7.64.1"},
	}

	assert.True(t, MatchKeyToMultiValue(exp, headers))

	exp["User-Agent"] = []string{"chrome"}

	assert.False(t, MatchKeyToMultiValue(exp, headers))
	assert.False(t, MatchKeyToMultiValue(map[string][]string{"content-type": {"text/plain"}}, map[string][]string{"content-type": {"text/plain", "text/csv"}}))
}

func TestMatchKeyToValue(t *testing.T) {
	exp := model.KeyToValue{
		"domain":     "ya.ru",
		"session_id": "123",
	}

	cookies := model.KeyToValue{
		"domain":     "ya.ru",
		"token":      "2412515151251521",
		"session_id": "123",
	}

	assert.True(t, MatchKeyToValue(exp, cookies))

	exp["token"] = ""
	assert.False(t, MatchKeyToValue(exp, cookies), "all expected fields should match")
}
