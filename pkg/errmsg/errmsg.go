package errmsg

import "fmt"

const (
	NoErrorCode = 200

	ServiceErrCode = iota + 10000
	AuthErrorCode
	ParseErrorCode
	DuplicationErrorCode
	DatabaseErrorCode
	RedisErrorCode
	ConfigMissErrorCode

	FileReadErrorCode
	FilePathCreateErrorCode
	FileWriteErrorCode
	FileDeleteErrorCode

	OssUploadErrorCode

	VideoNotExistErrorCode
	CommentExistErrorCode
	UserNotExistErrorCode
	UserExistErrorCode

	MfaGenareteErrorCode
	CryptEncodeErrorCode
	MfaOptCodeErrorCode

	IsNotImageErrorCode
	IsNotVideoErrorCode

	WebsockChatWriteErrorCode
	WebsockChatReadErrorCode
	WebsockChatParseErrorCode
	WebsockUpgradeErrorCode
)

const (
	NoErrorMsg = "OK"

	ServiceErrorMsg     = "Service Error"
	AuthErrorMsg        = "Auth Error"
	ParseErrorMsg       = "Parse Error"
	DuplicationErrorMsg = "Duplication Request"
	DatabaseErrorMsg    = "Database Operation Error"
	RedisErrorMsg       = "Redis Operation Error"
	ConfigMissErrorMsg  = "Config Missing Error"

	FileReadErrorMsg       = "File Read Error"
	FilePathCreateErrorMsg = "Create File Path Error"
	FileWriteErrorMsg      = "File Write Error"
	FileDeleteErrorMsg     = "File Delete Error"

	OssUploadErrorMsg = "Oss Upload Error"

	VideoNotExistErrorMsg   = "Video Not Exist"
	CommentNotExistErrorMsg = "Comment Not Exist"
	UserNotExistErrorMsg    = "User Not Exist"
	UserExistErrorMsg       = "User Exist"

	MfaGenareteErrorMsg = "MFA Generate Error"
	CryptEncodeErrorMsg = "Crypt Encode Error"
	MfaOptCodeErrorMsg  = "MFA Opt Code Error"

	IsNotImageErrorMsg = "Please oss image"
	IsNotVideoErrorMsg = "Please oss video"

	WebsockChatWriteErrorMsg = "Websocket Chat Write Error"
	WebsockChatReadErrorMsg  = "Websocket Chat Read Error"
	WebsockChatParseErrorMsg = "Websocket Chat Parse Error"
	WebsockUpgradeErrorMsg   = "Websocket Upgrade Error"
)

type ErrorMessage struct {
	ErrorCode int64
	ErrorMsg  string
}

func (err ErrorMessage) Error() string {
	return fmt.Sprintf("%v, Code:%v", err.ErrorMsg, err.ErrorCode)
}

func (err ErrorMessage) WithMessage(msg string) ErrorMessage {
	return ErrorMessage{
		ErrorCode: err.ErrorCode,
		ErrorMsg:  msg,
	}
}

func NewErrorMessage(code int64, msg string) ErrorMessage {
	return ErrorMessage{ErrorCode: code, ErrorMsg: msg}
}

var (
	NoError = NewErrorMessage(NoErrorCode, NoErrorMsg)

	ServiceError     = NewErrorMessage(ServiceErrCode, ServiceErrorMsg)
	AuthError        = NewErrorMessage(AuthErrorCode, AuthErrorMsg)
	ParseError       = NewErrorMessage(ParseErrorCode, ParseErrorMsg)
	DuplicationError = NewErrorMessage(DuplicationErrorCode, DuplicationErrorMsg)
	DatabaseError    = NewErrorMessage(DatabaseErrorCode, DatabaseErrorMsg)
	RedisError       = NewErrorMessage(RedisErrorCode, RedisErrorMsg)
	ConfigMissError  = NewErrorMessage(ConfigMissErrorCode, ConfigMissErrorMsg)

	FileReadError       = NewErrorMessage(FileReadErrorCode, FileReadErrorMsg)
	FilePathCreateError = NewErrorMessage(FilePathCreateErrorCode, FilePathCreateErrorMsg)
	FileWriteError      = NewErrorMessage(FileWriteErrorCode, FileWriteErrorMsg)
	FileDeleteError     = NewErrorMessage(FileDeleteErrorCode, FileDeleteErrorMsg)

	OssUploadError = NewErrorMessage(OssUploadErrorCode, OssUploadErrorMsg)

	VideoNotExistError   = NewErrorMessage(VideoNotExistErrorCode, VideoNotExistErrorMsg)
	CommentNotExistError = NewErrorMessage(CommentExistErrorCode, CommentNotExistErrorMsg)
	UserNotExistError    = NewErrorMessage(UserNotExistErrorCode, UserNotExistErrorMsg)
	UserExistError       = NewErrorMessage(UserExistErrorCode, UserExistErrorMsg)

	MfaGenareteError = NewErrorMessage(MfaGenareteErrorCode, MfaGenareteErrorMsg)
	CryptEncodeError = NewErrorMessage(CryptEncodeErrorCode, CryptEncodeErrorMsg)
	MfaOptCodeError  = NewErrorMessage(MfaOptCodeErrorCode, MfaOptCodeErrorMsg)

	IsNotImageError = NewErrorMessage(IsNotImageErrorCode, IsNotImageErrorMsg)
	IsNotVideoError = NewErrorMessage(IsNotVideoErrorCode, IsNotVideoErrorMsg)

	WebsockChatWriteError = NewErrorMessage(WebsockChatWriteErrorCode, WebsockChatWriteErrorMsg)
	WeBsockChatReadError  = NewErrorMessage(WebsockChatReadErrorCode, WebsockChatReadErrorMsg)
	WebsockChatParseError = NewErrorMessage(WebsockChatParseErrorCode, WebsockChatParseErrorMsg)
	WebsockUpgradeError   = NewErrorMessage(WebsockUpgradeErrorCode, WebsockUpgradeErrorMsg)
)
