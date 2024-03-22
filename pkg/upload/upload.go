package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"work4/bootstrap/env"
)

func IsImage(data *multipart.FileHeader) error {
	file, _ := data.Open()
	buffer := make([]byte, 261)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}

	if filetype.IsImage(buffer) {
		return nil
	}
	return errors.New("请上传图片")
}

func IsVideo(data *multipart.FileHeader) error {
	file, _ := data.Open()
	buffer := make([]byte, 261) // 读取足够多的字节以便确定文件类型

	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	if filetype.IsVideo(buffer) {
		return nil
	}
	return errors.New("请上传图片视频")
}

func SaveFile(data *multipart.FileHeader, storePath, fileName string) (err error) {

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		// 路径不存在，创建路径
		err := os.MkdirAll(storePath, 0755)
		if err != nil {
			return fmt.Errorf("创建路径错误: %w", err)
		}
	}

	//打开本地文件
	dist, err := os.OpenFile(filepath.Join(storePath, fileName), os.O_RDWR|os.O_CREATE, 777)
	if err != nil {
		return fmt.Errorf("创建文件错误:%w", err)
	}
	defer func(dist *os.File) {
		_ = dist.Close()
	}(dist)

	src, err := data.Open()
	if err != nil {
		return fmt.Errorf("保存文件错误:%w", err)
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
		Scope: env.QiNiuBucket,
	}

	mac := auth.New(env.QiNiuAccessKey, env.QiNiuSecretKey)
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
		return "", err
	}

	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}

	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return "", fmt.Errorf("上传错误:%w", err)
	}

	err = os.Remove(localFile)
	if err != nil {
		return "", err
	}
	fmt.Println("File deleted successfully")

	return storage.MakePublicURL(env.QiNiuDomain, ret.Key), nil

}
