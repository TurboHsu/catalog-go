package qiniu

import (
	"catalog-go/config"
	"context"
	"log"
	"path/filepath"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func PutFile(path string) error {
	c := config.CONFIG.Store.Qiniu
	fileName := filepath.Base(path)

	key := c.UploadPath + "/" + fileName
	putPolicy := storage.PutPolicy{
		Scope: c.Bucket + ":" + key,
	}

	mac := auth.New(c.AccessKey, c.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = true

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "catalog pic",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, path, &putExtra)
	if err != nil {
		log.Printf("[E] Failed to upload file: %v\n", err)
		return err
	}
	log.Printf("[I] File uploaded to qiniu: %s\n", ret.Key)
	return nil
}

func RemoveFile(path string) error {
	c := config.CONFIG.Store.Qiniu
	mac := auth.New(c.AccessKey, c.SecretKey)
	cfg := storage.Config{
		UseHTTPS: true,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	key := c.UploadPath + "/" + path
	err := bucketManager.Delete(c.Bucket, key)
	if err != nil {
		log.Printf("[E] Failed to remove file: %v\n", err)
		return err
	}
	log.Printf("[I] File removed from qiniu: %s\n", key)
	return nil
}
