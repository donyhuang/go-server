package db

import (
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var connMap = make(map[string]clickhouse.Conn)

func InitClickhouse(c []Conf) error {
	for _, cConf := range c {
		conn, err := clickhouse.Open(&clickhouse.Options{
			Addr: strings.Split(cConf.Server, ","),
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
			Auth: clickhouse.Auth{
				Database: cConf.Database,
				Username: cConf.User,
				Password: cConf.Pass,
			},
			DialTimeout: 3 * time.Second,
		})
		if err != nil {
			return err
		}
		key := cConf.Name
		if key == "" {
			key = DefaultKey
		}
		connMap[key] = conn
	}
	return nil
}

func GetClickHouseConn() clickhouse.Conn {
	return GetClickHouseConnByName(DefaultKey)
}
func GetClickHouseConnByName(str string) clickhouse.Conn {
	return connMap[str]
}
func CloseClick() {
	for _, v := range connMap {
		v.Close()
	}
}
