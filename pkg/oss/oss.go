package oss

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
)

func IsImage(data *multipart.FileHeader) error {
	file, _ := data.Open()
	buffer := make([]byte, 261)
	_, err := file.Read(buffer)
	if err != nil {
		return errmsg.FileReadError
	}

	if filetype.IsImage(buffer) {
		return nil
	}

	return errmsg.IsNotImageError
}

func IsVideo(data *multipart.FileHeader) error {
	file, _ := data.Open()
	buffer := make([]byte, 261)

	_, err := file.Read(buffer)
	if err != nil {
		return errmsg.FileReadError
	}
	if filetype.IsVideo(buffer) {
		return nil
	}

	return errmsg.IsNotVideoError
}

func SaveFile(data *multipart.FileHeader, storePath, fileName string) (err error) {

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		// 路径不存在，创建路径
		err := os.MkdirAll(storePath, 0755)
		if err != nil {
			return errmsg.FilePathCreateError
		}
	}

	//打开本地文件
	dist, err := os.OpenFile(filepath.Join(storePath, fileName), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errmsg.FileWriteError
	}
	defer func(dist *os.File) {
		_ = dist.Close()
	}(dist)

	src, err := data.Open()
	if err != nil {
		return errmsg.FileWriteError.WithMessage(err.Error())
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)
	_, err = io.Copy(dist, src)

	return
}

func Upload(localFile, filename, userid, origin string) (string, error) {

	key := fmt.Sprintf("%s/%s/%s", origin, userid, filename)

	putPolicy := storage.PutPolicy{
		Scope: constants.QiNiuBucket,
	}

	mac := auth.New(constants.QiNiuAccessKey, constants.QiNiuSecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneXinjiapo
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}

	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		return "", errmsg.OssUploadError.WithMessage(err.Error())
	}

	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}

	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return "", errmsg.OssUploadError
	}

	err = os.Remove(localFile)
	if err != nil {
		return "", errmsg.FileDeleteError
	}

	return storage.MakePublicURL(constants.QiNiuDomain, ret.Key), nil

}
