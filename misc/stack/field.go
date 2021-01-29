package stack

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// 打印文件名 ： 行号
func String(depth int) string {
	pc, _, n, ok := runtime.Caller(depth + 1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v", filepath.Base(runtime.FuncForPC(pc).Name()), n)
}

// 结构化日志，key - value形式
type Fields map[string]interface{}

func Field(kvs ...interface{}) Fields {
	return make(Fields).Field(kvs...)
}

func (f Fields) Clone() Fields {
	n := make(Fields, len(f))
	for k, v := range f {
		n[k] = v
	}
	return n
}

func (f Fields) Merge(f2 Fields) {
	for k, v := range f2 {
		if _, ok := f[k]; !ok {
			f[k] = v
		}
	}
}

func (f Fields) Field(kvs ...interface{}) Fields {
	if len(kvs)&1 != 0 {
		panic("invalid kvs")
	}
	for i := 0; i < len(kvs); i += 2 {
		name := kvs[i].(string)
		f[name] = kvs[i+1]
	}
	return f
}
