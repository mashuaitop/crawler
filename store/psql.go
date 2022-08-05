package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	cfg := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logger.Info,
				Colorful:      true,
			}), //日志级别
		DisableAutomaticPing: false, //初始化完成后自动ping
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名
		},
		PrepareStmt: false, //取消SQL缓存
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  uri(),
		PreferSimpleProtocol: true, // 禁用隐式 prepared statement
	}), cfg)
	if err != nil {
		msg := fmt.Sprint("init postgres db error, errmsg: ", err.Error())
		panic(msg)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}

func uri() string {
	host := "192.168.1.88"
	port := 7432
	user := "classcool"
	password := "123456"
	db := "template1"
	sslmode := "disable"

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, db, sslmode)
}
