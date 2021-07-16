package main

import (
	"fmt"
	minio_open "github.com/minio/minio-go/v6"
	"github.com/zoujiepro/file-server/store/minio"
)

func main() {
	client := minio.GetClient()

	testListObjects(client)

}

func testUpload(client *minio_open.Client) {
	filePath := "D:\\tmp\\Go语言趣学指南.rar"

	//新上传一个对象
	_, err := client.FPutObject("test1", "/test/testfile1", filePath, minio_open.PutObjectOptions{
		ContentType: "application/rar",
	})

	if err != nil {
		fmt.Println(err)
	}
}

func testListObjects(minioClient *minio_open.Client) {
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	isRecursive := true
	objectCh := minioClient.ListObjects("test1", "/test", isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}
}
