package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	dbMap = make(map[string]*gorm.DB)
)

func fixedConf(c Conf) Conf {
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = 100
	}
	if c.MaxLife == "" {
		c.MaxLife = "1h"
	} else {
		if _, err := time.ParseDuration(c.MaxLife); err != nil {
			c.MaxLife = "1h"
		}
	}
	return c
}
func InitMysql(c []Conf, serverMode string) error {
	for _, cConf := range c {
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&loc=Local", cConf.User, cConf.Pass, cConf.Server, cConf.Database)
		gormConfig := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			SkipDefaultTransaction: true,
		}
		if serverMode == "release" {
			gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		} else {
			gormConfig.Logger = logger.Default.LogMode(logger.Info)
		}
		db, err := gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			return err
		}
		var sqlDB *sql.DB
		sqlDB, err = db.DB()
		if err != nil {
			return err
		}
		cConf = fixedConf(cConf)
		duration, _ := time.ParseDuration(cConf.MaxLife)
		sqlDB.SetMaxIdleConns(int(cConf.MaxIdle))
		sqlDB.SetMaxOpenConns(int(cConf.MaxOpen))
		sqlDB.SetConnMaxLifetime(duration)
		key := cConf.Name
		if key == "" {
			key = DefaultKey
		}
		dbMap[key] = db

	}
	return nil
}

func GetDb() *gorm.DB {
	return GetDbByName(DefaultKey)
}
func GetDbByName(name string) *gorm.DB {
	return dbMap[name]
}
func GetPageOffset(page int, pageSize int) int {
	var offset int
	offset = (page - 1) * pageSize
	return offset
}
func CloseMysql() {
	for _, v := range dbMap {
		db, err := v.DB()
		if err != nil {
			continue
		}
		db.Close()
	}
}
