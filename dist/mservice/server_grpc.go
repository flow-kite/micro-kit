package mservice

import (
	"sync"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpc *grpc.Server

	unaryInterceptor []grpc.UnaryServerInterceptor // 服务端拦截器
	traceOnce        sync.Once
}
