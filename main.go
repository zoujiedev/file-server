package main

import (
	"github.com/zoujiepro/file-server/handler"
	"net/http"
)

func main() {
	//文件上传
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.UpdateHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.DeleteHandler))
	http.HandleFunc("/file/upload/fast", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	//分块上传
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	//用户
	http.HandleFunc("/user/signup", handler.UserSignUp)
	http.HandleFunc("/user/signin", handler.UserSignIn)
	http.HandleFunc("/user/files", handler.HTTPInterceptor(handler.GetUserFileHandler))

	http.ListenAndServe(":8080", nil)
}
