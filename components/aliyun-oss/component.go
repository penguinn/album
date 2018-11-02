package aliyun_oss

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/penguinn/penguin/component/config"
	"log"
	"os"
	"strings"
)

var bucket *oss.Bucket

func init() {
	client, err := oss.New(config.GetString("aliyun.endpoint"), config.GetString("aliyun.accessKeyID"), config.GetString("aliyun.accessKeySecret"))
	if err != nil {
		log.Fatal(err)
	}
	bucket, err = client.Bucket(config.GetString("aliyun.bucket"))
	if err != nil {
		log.Fatal(err)
	}
}

func GetSignURL(name string) (string, error) {

	signedURL, err := bucket.SignURL(name, oss.HTTPPut, 60)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

func BytesUpload(name string, content []byte) error {
	err := bucket.PutObject(name, bytes.NewReader(content))
	return err
}

func StringUpload(name, content string) error {
	err := bucket.PutObject(name, strings.NewReader(content))
	return err
}

func FileUpload(name, file string) error {
	defer os.Remove(file)
	err := bucket.PutObjectFromFile(name, file)
	return err
}
