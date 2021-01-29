package mconsul

import "github.com/hashicorp/consul/api"

type Service struct {
	Name     string
	Method   []string
	LanAddr  string
	WanAddr  string
	NodeName string

	id   string
	node *api.CatalogNode
}

func NewService(name string, method []string) *Service {
	return &Service{Name: name, Method: method}
}
