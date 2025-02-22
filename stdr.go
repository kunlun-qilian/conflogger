package conflogger

import (
	"context"
	"fmt"

	"github.com/go-courier/logr"
	"github.com/sirupsen/logrus"
)

func StdLogger() logr.Logger {
	return &stdLogger{lvl: logr.DebugLevel}
}

type stdLogger struct {
	lvl          logr.Level
	spans        []string
	keyAndValues []interface{}
}

func (d *stdLogger) SetLevel(lvl logr.Level) {
	d.lvl = lvl
}

func (d *stdLogger) WithValues(keyAndValues ...interface{}) logr.Logger {
	return &stdLogger{lvl: d.lvl, spans: d.spans, keyAndValues: append(d.keyAndValues, keyAndValues...)}
}

func (d *stdLogger) Start(ctx context.Context, name string, keyAndValues ...interface{}) (context.Context, logr.Logger) {
	return ctx, &stdLogger{lvl: d.lvl, spans: append(d.spans, name), keyAndValues: append(d.keyAndValues, keyAndValues...)}
}

func (d *stdLogger) End() {
	if len(d.spans) != 0 {
		d.spans = d.spans[0 : len(d.spans)-1]
	}
}

func (d *stdLogger) Debug(format string, args ...interface{}) {
	if logr.DebugLevel > d.lvl {
		return
	}
	// logrus.Debug(append(keyValues(append(d.keyAndValues, "level", "debug")), fmt.Sprintf(format, args...))...)
	logrus.Debug(fmt.Sprintf(format, args...))
}

func (d *stdLogger) Info(format string, args ...interface{}) {
	if logr.InfoLevel > d.lvl {
		return
	}

	logrus.Info(append(keyValues(append(d.keyAndValues, "level", "info")), fmt.Sprintf(format, args...))...)
}

func (d *stdLogger) Warn(err error) {
	if logr.WarnLevel > d.lvl {
		return
	}
	logrus.Warn(append(keyValues(append(d.keyAndValues, "level", "warn")), fmt.Sprintf("%v", err))...)
}

func (d *stdLogger) Error(err error) {
	if logr.ErrorLevel > d.lvl {
		return
	}
	logrus.Error(append(keyValues(append(d.keyAndValues, "level", "error")), fmt.Sprintf("%+v", err))...)
}

func keyValues(keyAndValues []interface{}) (values []interface{}) {
	if len(keyAndValues)%2 != 0 {
		return
	}
	for i := 0; i < len(keyAndValues); i += 2 {
		values = append(values, fmt.Sprintf("%v=%v", keyAndValues[i], keyAndValues[i+1]))
	}

	return
}
