package mconsul

import (
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/o-kit/micro-kit/misc/context"
)

// consul 数据中心 - 通过wan进行访问

// 默认为本地数据中心 - 一般局域网内只会访问本地的consul
var DefaultDatacenter = NewDatacenter(discoverConfig.Datacenter)

type Datacenter struct {
	Name string // 数据中心名称，第一个数据中心默认：dc1

	instance     *api.Client // 请求这个数据中心的客户端
	instanceOnce sync.Once   // 保证只能实例化一个
}

func NewDatacenter(name string) *Datacenter {
	return &Datacenter{
		Name: name,
	}
}

func GetDefaultDatacenter() *Datacenter {
	return DefaultDatacenter
}

// 获取 consul client
func (dc *Datacenter) getConsul() *api.Client {
	dc.instanceOnce.Do(func() {
		cli, err := api.NewClient(&api.Config{
			Address:    discoverConfig.Address,
			Datacenter: dc.Name,
		})
		if err != nil {
			panic(err)
		}
		dc.instance = cli
	})
	return dc.instance
}

// 获取数据中心的写权限 - 添加token验证
func (dc *Datacenter) getWriteOption() *api.WriteOptions {
	return &api.WriteOptions{
		Token:      discoverConfig.Token,
		Datacenter: dc.Name,
	}
}

// 获取数据中心的读取权限 - 添加token验证
func (dc *Datacenter) getQueryOption() *api.QueryOptions {
	return &api.QueryOptions{
		Token:      discoverConfig.Token,
		Datacenter: dc.Name,
	}
}

func GetDatacenters() ([]string, error) {
	return DefaultDatacenter.GetDatacenters()
}

// https://www.consul.io/api-docs/catalog#list-datacenters
// 获取多数据中心
func (dc *Datacenter) GetDatacenters() ([]string, error) {
	start := time.Now()
retry:
	ret, err := dc.getConsul().Catalog().Datacenters()
	if ok := dc.handleError(start, err); ok {
		goto retry
	}
	if len(ret) > 1 {
		// 可以对多个数据中心进行排序，这将影响数据中心的顺序
	}
	return ret, err
}

func (dc *Datacenter) GetNode() (*api.CatalogNode, error) {
	client := dc.getConsul()
	name, err := client.Agent().NodeName()
	if err != nil {
		return nil, err
	}
	node, _, err := client.Catalog().Node(name, dc.getQueryOption())
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("unable to fetch node: " + name)
	}
	return node, nil
}

func (dc *Datacenter) handleError(start time.Time, err error) bool {
	if err == nil {
		return false
	}
	if time.Since(start) > 10*time.Second {
		return false
	}

	if strings.Contains(err.Error(), "connection refused") {
		context.LogInfof("connect consul, retry, %v", err)
		time.Sleep(time.Second)
		return true
	}
	if strings.Contains(err.Error(), "No cluster leader") {
		context.LogInfof("no cluster leader, retry, %v", err)
		time.Sleep(time.Second)
		return true
	}
	if strings.Contains(err.Error(), "unable to fetch node") {
		context.LogInfof("no cluster leader, %v", err)
		time.Sleep(time.Second)
		return true
	}
	return false
}

func GetNodeFromDC(dc *Datacenter) (*api.CatalogNode, error) {
	start := time.Now()

retry:
	node, err := dc.GetNode()
	if ok := dc.handleError(start, err); ok {
		goto retry
	}
	return node, nil
}
