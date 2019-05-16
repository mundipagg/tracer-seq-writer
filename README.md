# Tracer Seq Writer

Custom writer for the [Tracer](https://github.com/mralves/tracer) log library


# Installation

* Using [go mod](https://github.com/golang/go/wiki/Modules) (*RECOMMENDED*):

Just edit your *go.mod* like the example below 
```
module github.com/user/repository
go 1.12

require (
    github.com/mundipagg/tracer-seq-writer version // ADD THIS LINE
)
```

* Using [dep](https://github.com/golang/dep):
```
$ dep ensure -add github.com/mundipagg/tracer-seq-writer@version
```

* Using [go get](https://golang.org/doc/articles/go_command.html#tmp_3) (*NOT RECOMMENDED*):
```
$ go get https://github.com/mundipagg/tracer-seq-writer
```


# Usage

## Configuration

|Field|Type|Mandatory?|Default|Description|
|---|---|---|:---:|---|
|Address|string|Y||Seq **full** endpoint (i.e. http://seq.io/api/events/eaw)|
|Key|string|N|""|Seq [API Key](https://docs.datalust.co/docs/api-keys)|
|Application|string|Y||Application name|
|Minimum Level|uint8|N|DEBUG|Minimum Level to log following the [syslog](https://en.wikipedia.org/wiki/Syslog#Severity_level) standard|
|Timeout|time.Duration|N|0 (infinite)|Timeout of the HTTP client|
|MessageEnvelop|string|N|"%v"|A envelop that *wraps* the original message|
|DefaultProperties|seq.Entry|N|{}|A generic object to append to *every* log entry, but can be overwritten by the original log entry|
|Buffer.Cap|int|N|100| Maximum capacity of the log buffer, when the buffer is full all logs are sent at once|
|Buffer.OnWait|int|N|100| Maximum size of the queue to send to Seq|
|Buffer.BackOff|time.Duration|N|60 seconds| Delay between retries to Seq|

## Examples

```go
package main

import (
	"github.com/mralves/tracer"
	"github.com/mundipagg/tracer-seq-writer"
)

func main() {
	logger := tracer.GetLogger("logger")
	logger.RegisterWriter(seq.New(seq.Config{
		Buffer:  seq.BufferConfig{Cap: 500, OnWait: 100, BackOff: 10000, Expiration: 10000},
		Address: "http://localhost:5341/api/events/raw",
	}))
	logger.Info("my beautiful log")
}
```