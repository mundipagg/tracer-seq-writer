// +build unit

package buffer

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestBuffer_Write(t *testing.T) {
	t.Parallel()
	t.Run("when the buffer is not full", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		subject := &buffer{
			Locker: &sync.Mutex{},
			size:   0,
			cap:    10,
			items:  make([]interface{}, 10),
			chunks: make(chan []interface{}, 10),
		}
		subject.Write("something")
		is.Equal(1, subject.size, "it should increment the size of the buffer in one unit")
		is.Equal([]interface{}{"something", nil, nil, nil, nil, nil, nil, nil, nil, nil}, subject.items, "it should change the buffer's inner slice")
	})
	t.Run("when the buffer is full", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		subject := &buffer{
			Locker: &sync.Mutex{},
			size:   0,
			cap:    1,
			items:  make([]interface{}, 1),
			chunks: make(chan []interface{}, 10),
		}
		subject.Write("something")
		is.Equal(0, subject.size, "it should remain zero")
		is.Equal([]interface{}{nil}, subject.items, "it should clean the buffer's inner slice")
		timeout := time.NewTimer(10 * time.Millisecond)
		select {
		case actual := <-subject.chunks:
			is.Equal([]interface{}{"something"}, actual, "it should read the expected slice")
		case <-timeout.C:
			is.Fail("nothing was published")
		}
	})
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("when the buffer expires", func(t *testing.T) {
		t.Parallel()
		t.Run("but the consumer returns an error", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			called := make(chan []interface{}, 1)
			err := errors.New("")
			subject := New(Config{
				OnOverflow: func(items []interface{}) error {
					called <- items
					err1 := err
					err = nil
					return err1
				},
				BackOff:    10 * time.Millisecond,
				Expiration: 10 * time.Millisecond,
				Cap:        10,
				OnWait:     10,
			})
			subject.Write(1)
			subject.Write(2)
			subject.Write(3)
			timeout := time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
			timeout = time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
		})
		t.Run("and the consumer returns no error", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			called := make(chan []interface{}, 1)
			subject := New(Config{
				OnOverflow: func(items []interface{}) error {
					called <- items
					return nil
				},
				BackOff:    10 * time.Millisecond,
				Expiration: 10 * time.Millisecond,
				Cap:        10,
				OnWait:     10,
			})
			subject.Write(1)
			subject.Write(2)
			subject.Write(3)
			timeout := time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
		})

	})
	t.Run("when the buffer overflow", func(t *testing.T) {
		t.Run("but the consumer returns an error", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			called := make(chan []interface{}, 1)
			err := errors.New("")
			subject := New(Config{
				OnOverflow: func(items []interface{}) error {
					called <- items
					err1 := err
					err = nil
					return err1
				},
				BackOff:    10 * time.Millisecond,
				Expiration: 100 * time.Millisecond,
				Cap:        3,
				OnWait:     10,
			})
			subject.Write(1)
			subject.Write(2)
			subject.Write(3)
			timeout := time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
			timeout = time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
		})
		t.Run("and the consumer returns no error", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			called := make(chan []interface{}, 1)
			subject := New(Config{
				OnOverflow: func(items []interface{}) error {
					called <- items
					return nil
				},
				BackOff:    10 * time.Millisecond,
				Expiration: 100 * time.Millisecond,
				Cap:        3,
				OnWait:     10,
			})
			subject.Write(1)
			subject.Write(2)
			subject.Write(3)
			timeout := time.NewTimer(20 * time.Millisecond)
			select {
			case items := <-called:
				is.Equal([]interface{}{1, 2, 3}, items, "it should return the expected array")
			case <-timeout.C:
				is.Fail("nothing was published")
			}
		})
	})
}
