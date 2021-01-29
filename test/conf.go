package test

import "github.com/o-kit/micro-kit/dist/mconf"

var instance Config

func Get() *Config {
	mconf.ByFlagOnce(&instance)

	return &instance
}

type Config struct {
	DBs           map[string]*DBConfig `toml:"dbs"`
	Redis         []string             `toml:"redis"`
	Elasticsearch Elasticsearch        `toml:"elasticsearch"`
}

type DBConfig struct {
	DriverName  string `toml:"driverName"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	UserName    string `toml:"username"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	MaxOpenConn int    `toml:"maxOpenConn"`
}

// Elasticsearch config for es
type Elasticsearch struct {
	URLs     []string `toml:"url"`
	Index    string   `toml:"index"`
	Type     string   `toml:"type"`
	UserName string   `toml:"username"`
	Password string   `toml:"password"`
}
