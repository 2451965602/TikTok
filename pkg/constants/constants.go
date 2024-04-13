package constants

var (
	MySQLUserName string
	MySQLPassWord string
	MySQLHost     string
	MySQLPort     string
	MySQLName     string
	MySQLDSN      string

	RedisUserName string
	RedisPassWord string
	RedisHost     string
	RedisPort     string

	QiNiuBucket    string
	QiNiuAccessKey string
	QiNiuSecretKey string
	QiNiuDomain    string
)

const (
	UserTable    = "user"
	VideoTable   = "video"
	CommentTable = "comment"
	LikeTable    = "like"
	SocialTable  = "social"
	MsgTable     = "messages"

	ContextUid  = "userid"
	DefaultTime = "1970-01-01 08:00:00"
	TimeFormat  = "2006-01-02 15:04:05"
)
