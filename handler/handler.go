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
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("error recv form file: %s", err.Error())
		utils.WriteFail(w, "error recv form file"+err.Error())
		return
	}

	fileStorePath := "D:\\tmp\\" + header.Filename
	metaInfo := meta.FileMeta{
		FileName:   header.Filename,
		Location:   fileStorePath,
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileStorePath)
	if err != nil {
		fmt.Printf("create tmp file fail : %s", err.Error())
		utils.WriteFail(w, "create tmp file fail : "+err.Error())
		return
	}

	defer newFile.Close()

	metaInfo.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("copy file fail: %s", err.Error())
		utils.WriteFail(w, "copy file fail: "+err.Error())
		return
	}

	metaInfo.FileSha1, err = utils.SHA1File(fileStorePath)
	//meta.UpdateFileMeta(metaInfo)
	meta.UpdateFileMetaDB(metaInfo)

	fmt.Printf("upload success file[%s], the sha1 = %s", metaInfo.FileName, metaInfo.FileSha1)
	utils.WriteSuccess(w, metaInfo)
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["filehash"][0]

	//fileMeta := meta.GetFileMeta(fileHash)
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		utils.WriteFail(w, "there are not exist meta info whit sha1 = "+fileHash)
		return
	}

	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Printf("json fail : %s", err.Error())
		utils.WriteFail(w, "internal server error: "+err.Error())
		return
	}

	utils.WriteSuccess(w, string(marshal))
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sha1 := r.FormValue("filehash")

	fileMeta, err := meta.GetFileMetaDB(sha1)
	if err != nil {
		fmt.Printf("query meta info err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	currentFileMeta, err := meta.GetFileMetaDB(sha1)
	if err != nil {
		fmt.Printf("query metainfo err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	currentFileMeta.FileName = newName
	if success := meta.UpdateFileMetaDB(currentFileMeta); !success {
		utils.WriteFail(w, "update fail")
		return
	}

	marshal, err := json.Marshal(currentFileMeta)
	if err != nil {
		fmt.Printf("json fail: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteSuccess(w, marshal)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	sha1 := r.Form.Get("filehash")
	_, err := meta.GetFileMetaDB(sha1)
	if err != nil {
		fmt.Println("query err: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if success := meta.DeleteFileMetaDB(sha1); success {
		utils.WriteSuccess(w, nil)
	} else {
		utils.WriteFail(w, "delete fail")
	}
}
