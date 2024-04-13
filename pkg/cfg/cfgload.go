package cfg

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"work4/biz/dal"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
)

var Config *viper.Viper

func Init() error {
	Config = viper.New()
	Config.SetConfigFile("./config/config.json")
	Config.SetConfigType("json")
	if err := Config.ReadInConfig(); err != nil {
		return errmsg.ConfigMissError
	}
	loadConfig()

	go func() {
		Config.WatchConfig()
		Config.OnConfigChange(func(e fsnotify.Event) {
			loadConfig()
		})
	}()

	return nil
}

func loadConfig() {

	constants.MySQLUserName = Config.GetString("MySQL.UserName")
	constants.MySQLPassWord = Config.GetString("MySQL.PassWord")
	constants.MySQLHost = Config.GetString("MySQL.Host")
	constants.MySQLPort = Config.GetString("MySQL.Port")
	constants.MySQLName = Config.GetString("MySQL.Name")
	constants.MySQLDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", constants.MySQLUserName, constants.MySQLPassWord, constants.MySQLHost, constants.MySQLPort, constants.MySQLName)
	dal.MysqlInit()

	constants.RedisUserName = Config.GetString("Redis.UserName")
	constants.RedisPassWord = Config.GetString("Redis.PassWord")
	constants.RedisHost = Config.GetString("Redis.Host")
	constants.RedisPort = Config.GetString("Redis.Port")
	dal.RedisInit()

	constants.QiNiuBucket = Config.GetString("QiNiu.Bucket")
	constants.QiNiuAccessKey = Config.GetString("QiNiu.AccessKey")
	constants.QiNiuSecretKey = Config.GetString("QiNiu.SecretKey")
	constants.QiNiuDomain = Config.GetString("QiNiu.Domain")
}
