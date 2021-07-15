package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

//获取s3连接
func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	//1. 初始化ceph的信息
	//通过查看 radosgw-admin user create --uid=user1 --display
	auth := aws.Auth{
		AccessKey: "xxx",
		SecretKey: "xxx",
	}
	//2. 创建s3连接
	region := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http:xxx.xxx.xxx.xxx:xxx",
		S3Endpoint:           "http:xxx.xxx.xxx.xxx:xxx",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	cephConn = s3.New(auth, region)
	return cephConn
}

//获取Bucket
func GetCephBucket(bucket string) *s3.Bucket {
	connection := GetCephConnection()
	return connection.Bucket(bucket)
}
