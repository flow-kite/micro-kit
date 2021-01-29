package mconsul

import "os"

type Config struct {
	Datacenter string
	Address    string
	Token      string
}

var discoverConfig = Config{
	// Address:    os.Getenv("CONSUL"),
	// Datacenter: os.Getenv("CONSUL_DC"),
	// Token: os.Getenv("CONSUL_TOKEN"),
	Address:    "47.96.100.46:8500",
	Datacenter: "dc1",
	Token:      os.Getenv("CONSUL_TOKEN"),
}

func SetAddress(addr string) {
	discoverConfig.Address = addr
}

func GetAddress() string {
	return discoverConfig.Address
}

func GetDatacenter() string {
	return DefaultDatacenter.Name
}
