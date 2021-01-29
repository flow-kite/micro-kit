package log

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/o-kit/micro-kit/misc/stack"
)

// todo 需要添加日志分割功能，使用Hook

var std = NewEx(1)

type Entry struct {
	*logrus.Entry
	depth int
}

func New() *Entry {
	return NewEx(0)
}

func NewEx(depth int) *Entry {
	logger := logrus.New()

	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
	}

	return &Entry{
		Entry: logrus.NewEntry(logger),
		depth: depth,
	}
}

func (e *Entry) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		e.Entry.Logger.Level = logrus.DebugLevel
	case "info":
		e.Entry.Logger.Level = logrus.InfoLevel
	case "warn":
		e.Entry.Logger.Level = logrus.WarnLevel
	case "error":
		e.Entry.Logger.Level = logrus.ErrorLevel
	case "fatal":
		e.Entry.Logger.Level = logrus.FatalLevel
	}
}

func (e *Entry) Printf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Printf(format, args...)
}

func (e *Entry) Println(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth)).Println(args...)
}

func (e *Entry) Debug(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Debug(args...)
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Debugf(format, args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Info(args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Infof(format, args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Warn(args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Warnf(format, args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Error(args...)
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Errorf(format, args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Fatal(args...)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Fatalf(format, args...)
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{e.Entry.WithField(key, value), e.depth}
}

func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	return &Entry{e.Entry.WithFields(fields), e.depth}
}

// -------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------

func SetLevel(level string) {
	std.SetLevel(level)
}

func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

func Println(args ...interface{}) {
	std.Println(args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Infof(layout string, args ...interface{}) {
	std.Infof(layout, args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}
