package yaml

import (
	"io/ioutil"
	"strings"

	"github.com/chennqqi/goutils/consul"
	"github.com/chennqqi/goutils/utils"
	"gopkg.in/yaml.v2"

	"github.com/chennqqi/go-ucfg"
)

func NewConfig(in []byte, opts ...ucfg.Option) (*ucfg.Config, error) {
	var m interface{}
	if err := yaml.Unmarshal(in, &m); err != nil {
		return nil, err
	}

	return ucfg.NewFrom(m, opts...)
}

func NewConfigWithFile(name string, opts ...ucfg.Option) (*ucfg.Config, error) {
	input, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	opts = append([]ucfg.Option{ucfg.MetaData(ucfg.Meta{name})}, opts...)
	return NewConfig(input, opts...)
}

func NewConfigWithConsulFile(name, consulagent string, opts ...ucfg.Option) (*ucfg.Config, error) {
	appName := utils.ApplicationName()
	extOffset := strings.LastIndexByte(name, '.')
	netName := name
	if extOffset > 0 {
		if name[:extOffset] == appName {
			netName = "config/" + name
		}
	}
	consulapi := consul.NewConsulOp(consulagent)
	consulapi.Fix()
	err := consulapi.Ping()

	var input []byte
	if err == nil {
		input, err = consulapi.Get(netName)
	} else {
		input, err = ioutil.ReadFile(name)
	}
	if err != nil {
		return nil, err
	}

	opts = append([]ucfg.Option{ucfg.MetaData(ucfg.Meta{name})}, opts...)
	return NewConfig(input, opts...)
}
