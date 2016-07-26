// +build windows

package logruswindows

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/svc/eventlog"
)

// EventHook sends logs to window event logs
type EventHook struct {
	*eventlog.Log
	source string
	levels []logrus.Level
}

// NewEventHook creates an event logging hook from even source
// and supported log levels
func NewEventHook(source string, levels []logrus.Level) (*EventHook, error) {
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	if err := eventlog.InstallAsEventCreate(source, supports); err != nil {
		return nil, errors.Wrapf(err, "eventlog.InstallAsEventCreate source=%s", source)
	}

	l, err := eventlog.Open(source)
	if err != nil {
		return nil, errors.Wrapf(err, "eventlog.Open source=%s", source)
	}
	return &EventHook{
		Log:    l,
		source: source,
		levels: levels,
	}, nil
}

// Fire extracts logrus entry and sends to window event log
func (hook *EventHook) Fire(entry *logrus.Entry) error {
	msg, err := entry.String()
	var eventID uint32 = 1
	id, ok := entry.Data["event_id"].(string)
	if ok {
		// attempt to convert to uint32 type event id
		var id64 uint64
		id64, err = strconv.ParseUint(id, 10, 64)
		if err == nil {
			eventID = uint32(id64)
		}
	}

	switch entry.Level {
	case logrus.DebugLevel, logrus.InfoLevel:
		return errors.Wrap(hook.Info(eventID, msg), "hook.Info")
	case logrus.WarnLevel:
		return errors.Wrap(hook.Warning(eventID, msg), "hook.Warning")
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return errors.Wrap(hook.Error(eventID, msg), "hook.Error")
	}
	return errors.Errorf("Unknown logrus level %s", entry.Level.String())
}

// Levels returns current available logging levels.
func (hook *EventHook) Levels() []logrus.Level {
	return hook.levels
}

// Close event log & removes registry
func (hook *EventHook) Close() error {
	if err := hook.Log.Close(); err != nil {
		return errors.Wrapf(err, "eventlog.Log.Close()")
	}
	return errors.Wrap(eventlog.Remove(hook.source), "eventlog.Remove")
}
