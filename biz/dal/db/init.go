package db

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok/pkg/constants"
	"tiktok/pkg/errmsg"
	"time"
)

var DB *gorm.DB

func Init() error {
	var err error

	DB, err = gorm.Open(mysql.Open(constants.MySQLDSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return errmsg.MysqlInitError.WithMessage(err.Error())
	}

	sqlDB, err := DB.DB()

	if err != nil {
		return errmsg.MysqlInitError.WithMessage(err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	if err = sqlDB.Ping(); err != nil {
		hlog.Error("数据库Ping失败: ", err)
		return errmsg.MysqlInitError.WithMessage(err.Error()) // 返回自定义错误消息
	}
	hlog.Info("MsSQL连接成功")

	return nil
}
