package context

import (
	"context"
	"sync"
)

func newctx(c context.Context) T {
	return T{Context: c, m: new(sync.Mutex)}
}

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
