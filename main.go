package main

import (
	"filestore-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadFileSuc)

	http.ListenAndServe(":8080", nil)

}
