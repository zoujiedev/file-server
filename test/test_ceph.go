package main

import (
	"fmt"
	"github.com/zoujiepro/file-server/store/ceph"
	"gopkg.in/amz.v1/s3"
)

func main() {
	bucket := ceph.GetCephBucket("testBucket1")

	//创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		fmt.Printf("create bucket[%s] err: %s\n", bucket, err.Error())
		return
	}

	//查询这个bucket下指定条件的object keys
	result, err := bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket[%s] list err: %s\n", bucket, err.Error())
		return
	}
	fmt.Printf("object keys: %+v\n", result)

	//新上传一个对象
	err = bucket.Put("/test/testfile1", []byte("this is a test file"), "octet-stream", s3.PublicRead)
	if err != nil {
		fmt.Printf("bucket[%s] put err: %s\n", bucket, err.Error())
		return
	}

	//查询这个bucket下指定条件的object keys
	result, err = bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket[%s] list err: %s\n", bucket, err.Error())
		return
	}
	fmt.Printf("object keys: %+v\n", result)
}
