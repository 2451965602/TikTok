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
	ContextUid   = "userid"
)
