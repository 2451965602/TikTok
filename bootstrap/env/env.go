package env

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	MySQLDSN string
	RedisDSN string

	UserTable    = "user"
	VideoTable   = "video"
	CommentTable = "comment"
	LikeTable    = "like"
	SocialTable  = "social"
	ContextUid   = "userid"

	QiNiuBucket    string
	QiNiuAccessKey string
	QiNiuSecretKey string
	QiNiuDomain    string
)

func Init() {

	file, err := ini.Load("config/config.ini")
	if err != nil {
		panic(fmt.Errorf("配置文件读取错误，请检查文件路径:%w", err))
	}

	MySQLUserName := file.Section("MySQL").Key("UserName").String()
	MySQLPassWord := file.Section("MySQL").Key("PassWord").String()
	MySQLHost := file.Section("MySQL").Key("Host").String()
	MySQLPort := file.Section("MySQL").Key("Port").String()
	MySQLName := file.Section("MySQL").Key("Name").String()

	MySQLDSN = MySQLUserName + ":" + MySQLPassWord + "@tcp(" + MySQLHost + ":" + MySQLPort + ")/" + MySQLName + "?charset=utf8mb4&parseTime=true"

	RedisUserName := file.Section("Redis").Key("UserName").String()
	RedisPassWord := file.Section("Redis").Key("PassWord").String()
	RedisHost := file.Section("Redis").Key("Host").String()
	RedisPort := file.Section("Redis").Key("Port").String()
	RedisDB := file.Section("Redis").Key("DB").String()

	RedisDSN = "redis://" + RedisUserName + ":" + RedisPassWord + "@" + RedisHost + ":" + RedisPort + "/" + RedisDB

	QiNiuBucket = file.Section("QiNiu").Key("Bucket").String()
	QiNiuAccessKey = file.Section("QiNiu").Key("AccessKey").String()
	QiNiuSecretKey = file.Section("QiNiu").Key("SecretKey").String()
	QiNiuDomain = file.Section("QiNiu").Key("Domain").String()

}
