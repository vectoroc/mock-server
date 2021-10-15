package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name   string
		result bool
		a      []string
		b      []string
	}{
		{"one element slices", true, []string{"x"}, []string{"x"}},
		{"multiple element slices", true, []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"b is subset of a", false, []string{"a", "b", "c"}, []string{"a", "b"}},
		{"a is subset of a", false, []string{"b", "c"}, []string{"a", "b", "c"}},
		{"mismatch", false, []string{"a", "b"}, []string{"c", "d"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.result, Equal(test.a, test.b))
		})
	}
}

func TestJSONBytesEqual(t *testing.T) {
	t.Run("it should ignore space in json", func(t *testing.T) {
		equal, err := JSONBytesEqual([]byte(`{"x":  1}`), []byte(`{ "x" :1}`))
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("it should match nested structs", func(t *testing.T) {
		equal, err := JSONBytesEqual([]byte(`{"x": 1, "arr": [1, true], "bar": {"baz": null}}`), []byte(`{"x": 1, "arr": [1, true], "bar": {"baz": null}}`))
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("mismatched jsons", func(t *testing.T) {
		equal, err := JSONBytesEqual([]byte(`{"baz": null}`), []byte(`{"baz": true}`))
		assert.NoError(t, err)
		assert.False(t, equal)
	})

	t.Run("wrong json should throw parse error", func(t *testing.T) {
		_, err := JSONBytesEqual([]byte(`{baz: null}`), []byte(`{"baz": true}`))
		assert.Error(t, err)
	})
}
