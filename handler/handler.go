package handler

import (
	"encoding/json"
	"fmt"
	"github.com/zoujiepro/file-server/meta"
	"github.com/zoujiepro/file-server/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

		fileStorePath := "D:\\tmp\\" + header.Filename

		metaInfo := meta.FileMeta{
			FileName:   header.Filename,
			Location:   fileStorePath,
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileStorePath)
		if err != nil {
			fmt.Printf("create tmp file fial : %s", err.Error())
		}

		defer newFile.Close()

		metaInfo.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("copy file fail: %s", err.Error())
		}

		metaInfo.FileSha1, err = utils.SHA1File(fileStorePath)
		meta.UpdateFileMeta(metaInfo)

		fmt.Printf("upload success file[%s], the sha1 = %s", metaInfo.FileName, metaInfo.FileSha1)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadFileSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success and ")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["filehash"][0]

	fileMeta := meta.GetFileMeta(fileHash)

	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("json fail : %s", err.Error())
		io.WriteString(w, "internal server error: "+err.Error())
		return
	}

	io.WriteString(w, string(marshal))
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sha1 := r.FormValue("filehash")

	fileMeta := meta.GetFileMeta(sha1)

	downloadFile, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer downloadFile.Close()

	readAll, err := ioutil.ReadAll(downloadFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+"\""+fileMeta.FileName+"\"")
	w.Write(readAll)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	sha1 := r.Form.Get("filehash")
	newName := r.Form.Get("newname")

	currentFileMeta := meta.GetFileMeta(sha1)
	currentFileMeta.FileName = newName

	meta.UpdateFileMeta(currentFileMeta)

	marshal, err := json.Marshal(currentFileMeta)
	if err != nil {
		fmt.Printf("json fail: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshal)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	sha1 := r.Form.Get("filehash")
	currentFileMeta := meta.GetFileMeta(sha1)

	err := os.Remove(currentFileMeta.Location)
	if err != nil {
		fmt.Printf("remove file fail: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	meta.DeleteFileMeta(sha1)
	marshal, err := json.Marshal(currentFileMeta)
	if err != nil {
		fmt.Printf("json fail: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshal)
}
