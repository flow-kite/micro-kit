package mservice

type Server struct {
	serviceName string // 服务名称
	Port        int    // 端口号

	GrpcServer
	WebServer
}
