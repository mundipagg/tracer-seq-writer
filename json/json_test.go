// +build unit

package json

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewWithCaseStrategy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	called := 0
	subject := NewWithCaseStrategy(func(s string) string {
		called++
		return strings.ToUpper(s)
	})

	bytes, err := subject.Marshal(map[string]string{"a": "b"})
	is.Nil(err, "it should return no error")
	is.Equal(`{"A":"b"}`, string(bytes), "it should return the expected json")
	is.Equal(1, called, "it should call the strategy exactly one time")
}

func TestNew(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	subject := New()
	bytes, err := subject.Marshal(map[string]string{"a": "b"})
	is.Nil(err, "it should return no error")
	is.Equal(`{"a":"b"}`, string(bytes), "it should return the expected json")
}
