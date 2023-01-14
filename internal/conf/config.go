package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type (
	etcd struct {
		Scheme string   `yaml:"scheme"`
		Hosts  []string `yaml:"hosts"`
	}

	rpcSvc struct {
		UserRPC svc `yaml:"userRPC"`
	}

	svc struct {
		Hosts []string `yaml:"hosts"`
		Name  string   `yaml:"name"`
	}

	telegram struct {
		TokenID string `yaml:"tokenID"`
		GroupID string `yaml:"groupID"`
	}

	wss struct {
		Addr string `yaml:"addr"`
	}
)

var (
	_c struct {
		Etcd     etcd     `yaml:"etcd"`
		RPCSvc   rpcSvc   `yaml:"rpcSvc"`
		Telegram telegram `yaml:"telegram"`
		Wss      wss      `yaml:"wss"`
	}
)

func InitYaml() {
	data, err := os.ReadFile("./internal/conf/etcd-test.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &_c); err != nil {
		panic(err)
	}

	for _, h := range _c.Etcd.Hosts {
		fmt.Println(h)
	}

}

func Etcd() *etcd {
	return &_c.Etcd
}

func RPCSvc() *rpcSvc {
	return &_c.RPCSvc
}

func Telegram() *telegram {
	return &_c.Telegram
}

func Wss() *wss {
	return &_c.Wss
}
