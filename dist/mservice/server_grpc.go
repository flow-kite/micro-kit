package mservice

import (
	"net"
	"sync"

	c2 "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPC 提供的注册
type GrpcRegister interface {
	GetServer() *grpc.Server                       // 获取grpc's server
	RegisterService(name string, methods []string) // 注册服务及其提供的接口方法
}

// 注册一个服务包括：name + methods
type grpcServiceInfo struct {
	name    string
	methods []string
}

// GRPC 服务主体
type GrpcServer struct {
	// 一个GRPC可以提供注册多个service
	services []grpcServiceInfo

	grpc *grpc.Server

	unaryInterceptor []grpc.UnaryServerInterceptor // 服务端拦截器
	traceOnce        sync.Once
}

// 提供服务端拦截器 添加
func (s *GrpcServer) GrpcAddInterceptor(i grpc.UnaryServerInterceptor) {
	s.unaryInterceptor = append(s.unaryInterceptor, i)
}

// 往GRPC服务中添加待注册的service
func (s *GrpcServer) RegisterService(name string, methods []string) {
	s.services = append(s.services, grpcServiceInfo{name: name, methods: methods})
}

func (s *GrpcServer) GetServer() *grpc.Server {
	if s.grpc == nil {
		s.grpc = grpc.NewServer(grpc.UnaryInterceptor(s.serverInterceptor))
	}
	return s.grpc
}

// 启动GRPC服务
func (s *GrpcServer) Serve(ln net.Listener) {
	_ = s.GetServer().Serve(ln)
}

// 服务端拦截器 + 注入在GetServer时候
func (s *GrpcServer) serverInterceptor(ctx c2.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	return nil, nil
}
