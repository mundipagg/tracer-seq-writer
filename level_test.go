// +build unit

package seq

import (
	"github.com/mralves/tracer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevel(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	cases := map[uint8]string{
		tracer.Debug:         Debug,
		tracer.Informational: Information,
		tracer.Notice:        Warning,
		tracer.Warning:       Warning,
		tracer.Error:         Error,
		tracer.Critical:      Error,
		tracer.Alert:         Fatal,
		tracer.Fatal:         Fatal,
		9:                    Verbose,
	}
	for input, expected := range cases {
		actual := Level(input)
		is.Equal(expected, actual, "it should return the expected value")
	}
}
