package aliyun_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/penguinn/penguin/component/config"
)

func init() {

}

func GetSignURL(name string) (string, error) {
	client, err := oss.New(config.GetString("aliyun.endpoint"), config.GetString("aliyun.accessKeyID"), config.GetString("aliyun.accessKeySecret"))
    if err != nil {
        return "", err
    }
    bucket, err := client.Bucket(config.GetString("bucket"))
    if err != nil {
        return "", err
    }

    signedURL, err := bucket.SignURL(name, oss.HTTPPut, 60)
    if err != nil {
        return "", err
    }
    return signedURL, nil
}