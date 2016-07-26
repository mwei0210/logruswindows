# Logrus hook for event logging at windows

[![Build Status](https://travis-ci.org/meomap/logruswindows.svg?branch=master)](https://travis-ci.org/meomap/logruswindows)

## Usage
```go
import (
  "github.com/meomap/logruswindows"
  "github.com/Sirupsen/logrus"
)

func main() {
  log       := logrus.New()
  hook, err := logruswindow.NewEventHook(EVENT_SOURCE, []logrus.Level{
    logrus.PanicLevel,
    logrus.FatalLevel,
    logrus.ErrorLevel,
  })

  if err == nil {
    log.Hooks.Add(hook)
  }
  hook.Close()
}
```

## Run test
```
GOOS=windows GOARCH=386 go test -v .
```
