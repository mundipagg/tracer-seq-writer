//+build unit

package seq

import (
	"errors"
	"github.com/icrowley/fake"
	"github.com/jarcoal/httpmock"
	"github.com/mralves/tracer"
	"github.com/mundipagg/tracer-seq-writer/buffer"
	"github.com/mundipagg/tracer-seq-writer/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestWriter_Write(t *testing.T) {
	t.Parallel()
	t.Run("when the minimum level is higher than the log level received", func(t *testing.T) {
		t.Parallel()
		buf := &buffer.Mock{}
		ref := time.Now()
		stackTrace := tracer.GetStackTrace(3)
		subject := &Writer{
			buffer:       buf,
			minimumLevel: tracer.Error,
		}

		entry := tracer.Entry{
			Level:         tracer.Debug,
			Message:       "Message",
			StackTrace:    stackTrace,
			Time:          ref,
			Owner:         "owner",
			TransactionId: "Transaction",
			Args: []interface{}{
				"Arg",
				Entry{
					"Nested": "value",
				},
			},
		}
		subject.Write(entry)
		time.Sleep(30 * time.Millisecond)
		buf.AssertExpectations(t)
	})
	t.Run("when the minimum level is lower than the log level received", func(t *testing.T) {
		t.Parallel()
		buf := &buffer.Mock{}
		ref := time.Now()
		stackTrace := tracer.GetStackTrace(3)
		event := event{
			Level:           Error,
			MessageTemplate: "Before Message After",
			Properties: Entry{
				"string":     "Arg",
				"Nested":     "value",
				"Caller":     stackTrace[0].String(),
				"RequestKey": "Transaction",
				"Name":       "Default",
			},
			Timestamp: ref.UTC().Format(time.RFC3339Nano),
		}
		buf.On("Write", event).Return()
		subject := &Writer{
			buffer:         buf,
			minimumLevel:   tracer.Debug,
			messageEnvelop: "Before %v After",
			defaultProperties: Entry{
				"Name": "Default",
			},
		}

		entry := tracer.Entry{
			Level:         tracer.Critical,
			Message:       "Message",
			StackTrace:    stackTrace,
			Time:          ref,
			Owner:         "owner",
			TransactionId: "Transaction",
			Args: []interface{}{
				"Arg",
				Entry{
					"Nested": "value",
				},
			},
		}
		subject.Write(entry)
		time.Sleep(30 * time.Millisecond)
		buf.AssertExpectations(t)
	})
}

func TestWriter_Send(t *testing.T) {
	t.Parallel()
	t.Run("when there is an invalid field value in event", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		c := &http.Client{}
		httpmock.ActivateNonDefault(c)
		subject := &Writer{
			address:    "http://log.io/",
			client:     c,
			marshaller: json.New(),
		}
		err := subject.send([]interface{}{
			event{
				Properties: Entry{
					"C": make(chan int),
				},
			},
		})
		is.NotNil(err, "it should return an error")
	})
	t.Run("when the request fails", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		c := &http.Client{}
		httpmock.ActivateNonDefault(c)
		url := "http://log.io/" + fake.Password(8, 8, false, false, false)
		httpmock.RegisterResponder("POST", url, func(request *http.Request) (response *http.Response, err error) {
			return nil, errors.New("failed")
		})
		subject := &Writer{
			address:    url,
			client:     c,
			marshaller: json.New(),
		}
		err := subject.send([]interface{}{
			event{
				Properties: Entry{
					"C": 15,
				},
			},
		})
		is.NotNil(err, "it should return an error")
	})
	t.Run("when the request fails", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		c := &http.Client{}
		httpmock.ActivateNonDefault(c)
		url := "http://log.io/" + fake.Password(8, 8, false, false, false)
		httpmock.RegisterResponder("POST", url, func(request *http.Request) (response *http.Response, err error) {
			is.Equal(http.Header{
				"X-Seq-Apikey": []string{"key"},
				"Content-Type": []string{"application/json"},
			}, request.Header, "it should return the expected header")
			return nil, errors.New("failed")
		})
		subject := &Writer{
			address:    url,
			client:     c,
			marshaller: json.New(),
			key:        "key",
		}
		err := subject.send([]interface{}{
			event{
				Properties: Entry{
					"C": 15,
				},
			},
		})
		is.NotNil(err, "it should return an error")
	})
	t.Run("when the request return an status unexpected", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		c := &http.Client{}
		httpmock.ActivateNonDefault(c)
		url := "http://log.io/" + fake.Password(8, 8, false, false, false)
		httpmock.RegisterResponder("POST", url, func(request *http.Request) (response *http.Response, err error) {
			is.Equal(http.Header{
				"X-Seq-Apikey": []string{"key"},
				"Content-Type": []string{"application/json"},
			}, request.Header, "it should return the expected header")
			return httpmock.NewBytesResponse(502, nil), nil
		})
		subject := &Writer{
			address:    url,
			client:     c,
			marshaller: json.New(),
			key:        "key",
		}
		err := subject.send([]interface{}{
			event{
				Properties: Entry{
					"C": 15,
				},
			},
		})
		is.NotNil(err, "it should return an error")
	})
	t.Run("when the request return 201", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		c := &http.Client{}
		httpmock.ActivateNonDefault(c)
		url := "http://log.io/" + fake.Password(8, 8, false, false, false)
		httpmock.RegisterResponder("POST", url, func(request *http.Request) (response *http.Response, err error) {
			is.Equal(http.Header{
				"X-Seq-Apikey": []string{"key"},
				"Content-Type": []string{"application/json"},
			}, request.Header, "it should return the expected header")
			return httpmock.NewBytesResponse(201, nil), nil
		})
		subject := &Writer{
			address:    url,
			client:     c,
			marshaller: json.New(),
			key:        "key",
		}
		err := subject.send([]interface{}{
			event{
				Properties: Entry{
					"C": 15,
				},
			},
		})
		is.Nil(err, "it should return no error")
	})
}
