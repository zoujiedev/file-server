package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		file, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internal server error")
		}
		io.WriteString(w, string(file))
		return
	} else if "POST" == r.Method {
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("error recv form file: %s", err.Error())
		}

		newFile, err := os.Create("D:\\tmp\\" + header.Filename)
		if err != nil {
			fmt.Printf("create tmp file fial : %s", err.Error())
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("copy file fail: %s", err.Error())
		}

		http.Redirect(w, r, "/upload/suc", http.StatusFound)
	}
}

func UploadFileSuc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}
