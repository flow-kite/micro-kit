package common

import "github.com/golang/protobuf/proto"

func GenOption(s []byte) *ServiceOpDesc {
	var ret ServiceOpDesc
	proto.Unmarshal(s, &ret)
	return &ret
}
