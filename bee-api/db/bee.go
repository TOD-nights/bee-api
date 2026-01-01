package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	config2 "gitee.com/stuinfer/bee-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var instance *DB
var once sync.Once

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别：Silent、Error、Warn、Info
				Colorful:      true,        // 彩色打印
			},
		),
	})
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) GetDB() *gorm.DB {
	return d.db
}

func SetDB(db *gorm.DB) {
	instance = &DB{db: db}
}

func InitDB() error {
	var err error
	dbCfg := config2.AppConfigIns.DB
	instance, err = InitDBByConfig(dbCfg)
	if err != nil {
		return err
	}
	return nil
}

func InitDBByConfig(dbCfg *config2.AppDBConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
	return NewDB(dsn)
}

func GetDB() *gorm.DB {
	once.Do(func() {
		if instance != nil {
			return
		}
		if err := InitDB(); err != nil {
			panic(err)
		}
	})
	return instance.db
}
