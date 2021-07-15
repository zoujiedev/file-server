package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zoujiepro/file-server/config"
)

var ossClient *oss.Client

func OSSClient() *oss.Client {
	if ossClient != nil {
		return ossClient
	}

	ossClient, err := oss.New(config.OSSEndpoint, config.OSSAccessKey, config.OSSAccessKeySecret)
	if err != nil {
		fmt.Printf("oss new err: %s\n", err.Error())
		return nil
	}

	return ossClient
}

func Bucket() *oss.Bucket {
	client := OSSClient()
	if client != nil {
		bucket, err := client.Bucket(config.OSSBucket)
		if err != nil {
			fmt.Printf("bucket err: %s\n", err.Error())
			return nil
		}
		return bucket
	}
	return nil
}
