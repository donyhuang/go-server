package server

import (
	"github.com/donyhuang/go-server/pkg/prometheus"
	"github.com/donyhuang/go-server/pkg/server"
	"github.com/spf13/viper"

	"github.com/donyhuang/go-server/pkg/cache"
	"github.com/donyhuang/go-server/pkg/conf"
	"github.com/donyhuang/go-server/pkg/db"
	"github.com/donyhuang/go-server/pkg/log"
	kafka "github.com/donyhuang/go-server/pkg/mq"
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
