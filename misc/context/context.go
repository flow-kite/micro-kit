package context

import (
	"context"
	"fmt"
	"sync"

	"github.com/o-kit/micro-kit/misc/log"
	"github.com/o-kit/micro-kit/misc/stack"
)

func From(c context.Context) T {
	ctx := newctx(c)
	return ctx
}

func newctx(c context.Context) T {
	return T{Context: c, m: new(sync.Mutex)}
}

func Dump() T {
	ctx := newctx(context.TODO())
	return ctx
}

func WithCancel(parent T) (T, CancelFunc) {
	ctx, cancel := context.WithCancel(parent.Context)
	return newctx(ctx), CancelFunc(cancel)
}

type CancelFunc context.CancelFunc

type T struct {
	context.Context
	m *sync.Mutex
}

func (c *T) SetValue(key interface{}, value interface{}) {
	ctx := context.WithValue(c.Context, key, value)
	c.Context = ctx
}

func (c *T) GetValue(key interface{}) interface{} {
	return c.Context.Value(key)
}

var logger = log.NewEx(1)

func LogInfo(msg string, fields ...interface{}) {
	if len(fields)%2 != 0 {
		logger.Info(fmt.Sprintf(msg, fields...))
		return
	}
	logger.WithFields(stack.Field(fields...)).Info(msg)
}

func LogInfof(msg string, args ...interface{}) {
	logger.Info(fmt.Sprintf(msg, args...))
}
