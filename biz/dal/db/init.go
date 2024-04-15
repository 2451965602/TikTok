package db

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"work4/pkg/constants"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open(mysql.Open(constants.MySQLDSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	hlog.Info("MsSQL连接成功")

}
