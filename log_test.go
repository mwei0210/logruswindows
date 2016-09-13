// +build windows

package logruswindows

import (
	"io/ioutil"
	"testing"

	"github.com/Sirupsen/logrus"
)

const (
	src = "mylog"
	msg = "Errors happened!"
)

func TestEventHook(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	hook, err := NewEventHook(src, []logrus.Level{
		logrus.ErrorLevel,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if closeErr := hook.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}()

	logger.Hooks.Add(hook)

	fields := logrus.Fields{
		"func":   "DoSomething",
		"server": "localhost",
		"tag":    "development",
	}

	logger.WithFields(fields).Error(msg)

	t.Log("test with event id")
	fields["event_id"] = "321"
	logger.WithFields(fields).Error(msg)

	fields["event_id"] = "50b432ff-2be4-46c3-bfc3-df63ba299670"
	logger.WithFields(fields).Error(msg)

	fields["event_id"] = 432
	logger.WithFields(fields).Error(msg)

	_, err = NewEventHook(src, []logrus.Level{
		logrus.ErrorLevel,
	})
	if err != nil {
		t.Fatal(err)
	}
}
