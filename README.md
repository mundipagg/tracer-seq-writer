# SEQ Writer
This library is a writer(sink) SEQ to use with [Tracer](https://github.com/mralves/tracer)

## How to install
Using go get (not recommended):
```bash
go get github.com/mundipagg/tracer-seq-writer
```

Using [dep](github.com/golang/dep) (recommended):
```bash
dep ensure --add github.com/mundipagg/tracer-seq-writer@<version>
```

## How to use

Below follows a simple example of how to use this lib:

```go
package main

import (
	"fmt"
	"time"

	"github.com/mralves/tracer"
	seq "github.com/mundipagg/tracer-seq-writer"

	bsq "github.com/mundipagg/tracer-seq-writer/buffer"
)

type Safe struct {
	tracer.Writer
}

type LogEntry = map[string]interface{}

func (s *Safe) Write(entry tracer.Entry) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v", err)
		}
	}()
	s.Writer.Write(entry)
}

func configureTracer() {
	var writers []tracer.Writer
	tracer.DefaultContext.OverwriteChildren()
	writers = append(writers, &Safe{seq.New(seq.Config{
		Timeout:      3 * time.Second,
		MinimumLevel: tracer.Debug,
		DefaultProperties: LogEntry{
			"Application": "MyApplicationName",
			"Environment": "Production",
		},
		Key:         "yYka3WMkG5sfqVNlUUYj",
		Address:     "http://localhost:5341/api/events/raw",
		Buffer: bsq.Config{
			OnWait:     2,
			BackOff:    1 * time.Second,
			Expiration: 5 * time.Second,
		},
	})})

	for _, writer := range writers {
		tracer.RegisterWriter(writer)
	}
}

func inner() {
	logger := tracer.GetLogger("moduleA.inner")
	logger.Info("don't know which transaction is this")
	logger.Info("but this log in this transaction")
	logger = logger.Trace()
	go func() {
		logger.Info("this is also inside the same transaction")
		func() {
			logger := tracer.GetLogger("moduleA.inner.nested")
			logger.Info("but not this one...")

		}()
	}()
}

func main() {
	configureTracer()
	logger := tracer.GetLogger("moduleA")
	logger.Info("logging in transaction 'A'", "B")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'B'", "B")
	logger.Info("logging in transaction 'A'", "A")
	logger.Info("logging in transaction 'A'", "A")
	logger = logger.Trace("C") // now all logs on this logger will be on the transaction C
	logger.Info("logging in transaction 'C'")
	logger.Info("logging in transaction 'C'", "A")
	inner()
	fmt.Scanln()
}

```