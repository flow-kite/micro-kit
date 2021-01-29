package mservice

import (
	c2 "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"

	"github.com/o-kit/micro-kit/dist/proto/common"
)

// GRPC Client
type Client struct {
	GRPC *grpc.ClientConn
	cfg  *ClientConfig
}

type ClientConfig struct {
	Balancer balancer.Balancer
	Name     string
	Desc     *common.ServiceOpDesc // 这里拿到的是每个服务的 proto 文件中的配置信息
}

// 获取客户端默认配置
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{}
}

func (c *Client) clientInterceptor(ctx c2.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 这里是拦截器，可以做很多事情

	return nil
}

func (c *Client) Close() {
	c.GRPC.Close()
}

// 需要关注这个 name
func NewClientEx(name string, cfg *ClientConfig) (*Client, error) {
	if cfg == nil {
		cfg = new(ClientConfig)
	}
	cfg.Name = name
	cli := &Client{
		cfg: cfg,
	}

	var dialOption grpc.DialOption
	dialOption = grpc.WithUnaryInterceptor(cli.clientInterceptor)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		dialOption,
	}

	// 根据服务名称 - name 获取对应的 grpc client
	grpcClient, err := grpc.Dial(name, opts...)
	if err != nil {
		return nil, err
	}

	cli.GRPC = grpcClient

	return cli, nil
}
