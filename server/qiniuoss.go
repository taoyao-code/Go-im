package server

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
)

//var (
//	ossDomain = viper.GetString(`qiniu.QINIU_DOMAIN`)
//	accessKey = viper.GetString(`qiniu.QINIU_ACCESS_KEY`)
//	secretKey = viper.GetString(`qiniu.QINIU_SECRET_KEY`)
//	bucket    = viper.GetString(`qiniu.QINIU_TEST_BUCKET`)
//)

type UploadTokenService struct{}
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

//文档：https://developer.qiniu.com/kodo/sdk/1238/go#1
// UploadQiNiuYun 上传文件到七牛云
func (service *UploadTokenService) UploadQiNiuYun(filePath string, key string) (url string, err error) {

	ossDomain := viper.GetString(`qiniu.QINIU_DOMAIN`)
	accessKey := viper.GetString(`qiniu.QINIU_ACCESS_KEY`)
	secretKey := viper.GetString(`qiniu.QINIU_SECRET_KEY`)
	bucket := viper.GetString(`qiniu.QINIU_TEST_BUCKET`)

	// 简单上传凭证
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	// 设置上传凭证有效期
	putPolicy.Expires = 7200 //2小时有效期
	mac := auth.New(accessKey, secretKey)
	cfg := storage.Config{}
	//七牛云存储空间设置首页有存储区域
	cfg.Zone = &storage.ZoneHuadong
	//不启用HTTPS域名
	cfg.UseHTTPS = false
	//不使用CND加速
	cfg.UseCdnDomains = false
	//构建上传表单对象
	formUploader := storage.NewFormUploader(&cfg)
	//ret := storage.PutRet{}
	ret := MyPutRet{}
	// 可选
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "logo",
		},
	}
	upToken := putPolicy.UploadToken(mac)
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
	if err != nil {
		return "", err
	}
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, ossDomain, key, deadline)
	return privateAccessURL, nil
}
