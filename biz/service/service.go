package service

import (
	"github.com/cloudwego/hertz/pkg/app"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"
	"tiktok/pkg/constants"
	"tiktok/pkg/errmsg"

	"tiktok/pkg/oss"
)

func GetUidFormContext(c *app.RequestContext) int64 {
	uid, _ := c.Get(constants.ContextUid)
	userid, err := convertToInt64(uid)
	if err != nil {
		panic(err)
	}

	return userid
}

func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, errmsg.ParseError
	}
}

func UploadAndGetUrl(data *multipart.FileHeader, userid, sort string) (string, error) {
	ext := strings.ToLower(path.Ext(data.Filename))

	fileName := uuid.NewV4().String() + ext
	storePath := filepath.Join("static", userid, sort)

	if err := oss.SaveFile(data, storePath, fileName); err != nil {
		return "", err
	}

	url, err := oss.Upload(filepath.Join(storePath, fileName), fileName, userid, sort)

	if err != nil {
		return "", err
	}

	return url, nil
}
