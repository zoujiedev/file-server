package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zoujiepro/file-server/db"
	"github.com/zoujiepro/file-server/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	fileName := header.Filename

	newFile, err := os.Create(fileStorePath)
	if err != nil {
		fmt.Printf("create tmp file fail : %s", err.Error())
		utils.WriteFail(w, "create tmp file fail : "+err.Error())
		return
	}

	defer newFile.Close()

	fileSize, err := io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("copy file fail: %s", err.Error())
		utils.WriteFail(w, "copy file fail: "+err.Error())
		return
	}

	fileSha1, err := utils.SHA1File(fileStorePath)
	if uploadFile := db.UploadFile(fileSha1, fileName, fileSize, fileStorePath); !uploadFile {
		utils.WriteFail(w, "文件入库失败!")
		return
	}

	token := r.Header.Get("token")
	username, _ := db.GetUsernameByToken(token)

	if userFile := db.InsertUserFile(username, fileSha1, fileSize, fileName); !userFile {
		utils.WriteFail(w, "用户文件入库失败!")
		return
	}

	fmt.Printf("upload success file[%s], the sha1 = %s", fileName, fileSha1)
	utils.WriteSuccess(w, nil)
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["filehash"][0]

	tblf, err := db.GetFileMeta(fileHash)
	if err != nil {
		utils.WriteFail(w, "there are not exist meta info with sha1 = "+fileHash)
		return
	}

	marshal, err := json.Marshal(tblf)
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

	tblf, err := db.GetFileMeta(sha1)
	if err != nil {
		fmt.Printf("query meta info err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	downloadFile, err := os.Open(tblf.FileAddr.String)
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
	w.Header().Set("Content-Disposition", "attachment;filename="+"\""+tblf.FileName.String+"\"")
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

	tblf, err := db.GetFileMeta(sha1)
	if err != nil {
		fmt.Printf("query metainfo err: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if success := db.UpdateFileMeta(sha1, newName); !success {
		utils.WriteFail(w, "update fail")
		return
	}

	tblf.FileName = sql.NullString{
		newName, true,
	}

	marshal, err := json.Marshal(tblf)
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
	_, err := db.GetFileMeta(sha1)
	if err != nil {
		fmt.Println("query err: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if success := db.DeleteFileMeta(sha1); success {
		utils.WriteSuccess(w, nil)
	} else {
		utils.WriteFail(w, "delete fail")
	}
}

func GetUserFileHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	var username string
	var success bool
	if username, success = db.GetUsernameByToken(token); !success {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userFiles, success := db.GetUserFileByUserName(username)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteSuccess(w, userFiles)
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("token")
	username, _ := db.GetUsernameByToken(token)

	r.ParseForm()

	filehash := r.Form.Get("filehash")
	fileName := r.Form.Get("fileName")
	fileSize, err := strconv.ParseInt(r.Form.Get("fileSize"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		utils.WriteFail(w, "unable to query file with hash: "+filehash)
		return
	}

	if fileMeta.FileSize.Int64 != fileSize {
		utils.WriteFail(w, "file with hash: "+filehash+"size["+strconv.FormatInt(fileMeta.FileSize.Int64, 10)+"]not equal current parameter fileSize: "+strconv.FormatInt(fileSize, 10))
		return
	}

	fastUpload := db.InsertUserFile(username, filehash, fileSize, fileName)
	if !fastUpload {
		utils.WriteFail(w, "用户文件入库失败")
		return
	}
	utils.WriteSuccess(w, nil)
}
