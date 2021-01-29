package mconsul

import (
	"fmt"
	"testing"

	"github.com/treasure1993/base/micro/context"
)

func TestGetDatacenters(t *testing.T) {
	ret, err := GetDatacenters()
	if err != nil {
		context.LogInfof("err = %v", err)
	}

	fmt.Println(ret)
}
