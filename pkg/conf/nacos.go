package conf

import (
	"errors"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	client             config_client.IConfigClient
	confData           sync.Map
	NaCosConfNotExists = errors.New("naCos conf not exists")
)

type NacosConf struct {
	Server struct {
		Host string
		Port uint64
	}
	Client struct {
		TimeoutMs uint64 `yaml:"timeoutMs"`
		LogLevel  string `yaml:"logLevel"`
		LogDir    string `yaml:"logDir"`
		CacheDir  string `yaml:"cacheDir"`
		Namespace string `yaml:"namespace"`
	}
	Data []NacosItems
}

type NacosItems struct {
	Name string `yaml:"name"`
	Conf struct {
		DataId   string `yaml:"dataId"`
		Group    string `yaml:"group"`
		ConfType string `yaml:"confType"`
	} `yaml:"conf"`
}
type CancelListenFunc func() error

func InitNacos(c *NacosConf) ([]CancelListenFunc, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Server.Host, c.Server.Port),
	}
	var clientOption = []constant.ClientOption{
		constant.WithTimeoutMs(c.Client.TimeoutMs),
		constant.WithCacheDir(c.Client.CacheDir),
		constant.WithLogLevel(c.Client.LogLevel),
		constant.WithLogDir(c.Client.LogDir),
		constant.WithNamespaceId(c.Client.Namespace),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(clientOption...)
	var err error
	client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}
	return GetAndWatch(c.Data)
}

func GetAndWatch(c []NacosItems) ([]CancelListenFunc, error) {
	var sf []CancelListenFunc
	for _, ncItem := range c {
		nc := ncItem
		v, err := client.GetConfig(vo.ConfigParam{
			DataId: nc.Conf.DataId,
			Group:  nc.Conf.Group,
		})
		if err != nil {
			continue
		}
		codec := GetCodec(nc.Conf.ConfType)
		if codec == nil {
			continue
		}
		codeDataPointer := codec.GetEmptyPointer()
		err = codec.Unmarshal([]byte(v), codeDataPointer)
		if err != nil {
			return nil, err
		}
		confData.Store(nc.Name, codeDataPointer)
		param := vo.ConfigParam{
			DataId: nc.Conf.DataId,
			Group:  nc.Conf.Group,
			OnChange: func(namespace, group, dataId, data string) {
				codeDataPointer := codec.GetEmptyPointer()
				err := codec.Unmarshal([]byte(data), codeDataPointer)
				if err != nil {
					return
				}
				confData.Store(nc.Name, codeDataPointer)
			},
		}
		err = client.ListenConfig(param)
		if err != nil {
			return nil, err
		}
		sf = append(sf, func(param vo.ConfigParam) CancelListenFunc {
			return func() error {
				param.OnChange = nil
				return client.CancelListenConfig(param)
			}
		}(param))
	}
	return sf, nil
}

func GetNacosConf(name string) interface{} {
	v, ok := confData.Load(name)
	if !ok {
		return nil
	}
	return v
}

func PublishContent(dataId, group, content string) (bool, error) {
	confParam := vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	}
	return client.PublishConfig(confParam)
}
