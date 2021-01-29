package mconsul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// 注册 service 及其地址到 consul
func (s *Service) Register(port int, deregister bool, tags []string) error {
	node, err := GetNodeFromDC(DefaultDatacenter)
	if err != nil {
		return errors.New("failed to getNode: " + err.Error())
	}
	addr := node.Node.Address
	fmt.Println("address: ", addr)

	s.LanAddr = node.Node.TaggedAddresses["lan"]

	client := DefaultDatacenter.getConsul()

	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                s.id,   // 需要注册的服务ID example: oms:192.168.1.174:63333
		Name:              s.Name, // 需要注册的服务名称
		Address:           "",     // 服务IP地址
		Port:              port,   // 服务端口
		EnableTagOverride: true,
	})
}

// deregister 注销服务
func (s *Service) Deregister() error {
	if s.id == "" {
		return nil
	}
	return DefaultDatacenter.getConsul().Agent().ServiceDeregister(s.id)
}
