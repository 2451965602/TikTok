package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"work4/bootstrap/env"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open(mysql.Open(env.MySQLDSN), &gorm.Config{
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
}
