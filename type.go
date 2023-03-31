package server

import (
	"github.com/spf13/viper"
	"gitlab.nongchangshijie.com/go-base/server/pkg/prometheus"
	"gitlab.nongchangshijie.com/go-base/server/pkg/server"

	"gitlab.nongchangshijie.com/go-base/server/pkg/cache"
	"gitlab.nongchangshijie.com/go-base/server/pkg/conf"
	"gitlab.nongchangshijie.com/go-base/server/pkg/db"
	"gitlab.nongchangshijie.com/go-base/server/pkg/log"
	kafka "gitlab.nongchangshijie.com/go-base/server/pkg/mq"
)

var appConf AppConf

type AppConf struct {
	Server struct {
		Mode       server.Mode
		Service    []server.Service
		Prometheus prometheus.Config
	}
	Cache struct {
		Redis []cache.RedisConf
	}
	Log  log.Conf
	Conf struct {
		Nacos conf.NacosConf
	}
	Kafka   kafka.Config
	Storage Storage
}

type Storage struct {
	Clickhouse []db.Conf
	Mysql      []db.Conf
}

func GetGlobalConf() *AppConf {
	return &appConf
}
func InitConf() (*AppConf, error) {
	err := viper.Unmarshal(&appConf)
	return &appConf, err
}
