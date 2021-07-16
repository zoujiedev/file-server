package minio

import (
	"fmt"
	minio_oss "github.com/minio/minio-go/v6"
)

var client *minio_oss.Client
var endPoint string = "xxx.xxx.xxx.xxx:9000"
var accessKey string = "xxxx"
var accessKeySecret string = "xxxx"
var location string = "cn-northwest-1"

func GetClient() *minio_oss.Client {
	if client != nil {
		return client
	}

	c, err := minio_oss.New(endPoint, accessKey, accessKeySecret, false)
	if err != nil {
		fmt.Printf("minio connect err: %s\n", err.Error())
		return nil
	}
	client = c
	return client
}
