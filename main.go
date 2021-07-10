package main

import (
	"github.com/zoujiepro/file-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.UpdateHandler)
	http.HandleFunc("/file/delete", handler.DeleteHandler)

	http.HandleFunc("/user/signup", handler.UserSignUp)
	http.HandleFunc("/user/signin", handler.UserSignIn)

	http.ListenAndServe(":8080", nil)
}
