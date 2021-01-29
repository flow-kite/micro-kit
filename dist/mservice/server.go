package mservice

import (
	"fmt"
	"net"

	"github.com/soheilhy/cmux"

	"github.com/o-kit/micro-kit/dist/proto/common"
	"github.com/o-kit/micro-kit/misc/context"
)

type Options interface {
	GetServiceDesc(serviceName string) *common.ServiceOpDesc
	GetMethodDesc(service, method string) (*common.ServiceOpDesc, *common.MethodOpDesc)
	GetAllServiceDesc() []*common.ServiceOpDesc
}

type Server struct {
	options []*common.ServiceOpDesc

	serviceName string // 服务名称
	Port        int    // 端口号

	GrpcServer
	WebServer
}

// run GRPC and WebApi
func (s *Server) Run(ctx context.T) error {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
	if err != nil {
		return err
	}
	mux := cmux.New(ln)
	{
		ln := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		go s.GrpcServer.Serve(ln)
	}
	{
		ln := mux.Match(cmux.Any())
		go s.WebServer.Serve(ctx, ln)
	}

	return mux.Serve()
}

func (s *Server) GetAllServiceDesc() []*common.ServiceOpDesc {
	return s.options
}

func (s *Server) GetServiceDesc(service string) *common.ServiceOpDesc {
	for _, op := range s.options {
		if op.Name == service {
			return op
		}
	}
	return nil
}

func (s *Server) GetMethodDesc(serviceName, methodName string) (*common.ServiceOpDesc, *common.MethodOpDesc) {
	service := s.GetServiceDesc(serviceName)
	for _, method := range service.GetMethods() {
		if method.Name == methodName {
			return service, method
		}
	}
	return service, nil
}

// 通过调用Register将服务的Desc注入进来
func (s *Server) RegisterServiceDesc(op *common.ServiceOpDesc) {
	s.options = append(s.options, op)
	s.WebServer.options = s
	s.GrpcServer.options = s
}
