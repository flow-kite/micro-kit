package build

import "flag"

// 判断是不是测试环境
var IsTest = func() bool {
	return flag.Lookup("test.v") != nil
}()
