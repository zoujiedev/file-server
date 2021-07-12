package db

import (
	"fmt"
	dblayer "github.com/zoujiepro/file-server/db/mysql"
	"time"
)

type TableUserFile struct {
	UserName   string
	FileSha1   string
	FileSize   int64
	FileName   string
	UploadAt   time.Time
	LastUpdate time.Time
	status     int
}

func InsertUserFile(username string, fileHash string, fileSize int64, fileName string) bool {
	sqlStr := `insert into tbl_user_file (user_name,file_sha1,file_size,file_name)
				values (?,?,?,?)`

	fmt.Printf("sql: %s\n parameter: username = %s,filehash = %s, fileSize = %d, fileName = %s\n", sqlStr, username, fileHash, fileSize, fileName)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("sql prepare err: %s", err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(username, fileHash, fileSize, fileName)
	if err != nil {
		fmt.Printf("sql exec err: %s", err.Error())
		return false
	}

	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return true
	}
	fmt.Printf("sql exec success,but affected 0")
	return false
}

func GetUserFileByUserName(username string) ([]TableUserFile, bool) {
	sqlStr := `select user_name,file_sha1,file_size,file_name,upload_at,last_update,status 
               from tbl_user_file 
		  	   where user_name = ?`

	fmt.Printf("sql: %s\n parameter: username = %s\n", sqlStr, username)

	var result []TableUserFile
	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("sql prepare err: %s", err.Error())
		return result, false
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Printf("sql query err: %s", err.Error())
		return result, false
	}

	for rows.Next() {
		var tbl TableUserFile
		rows.Scan(&(tbl.UserName), &(tbl.FileSha1), &(tbl.FileSize), &(tbl.FileName), &(tbl.UploadAt), &(tbl.LastUpdate), &(tbl.status))
		result = append(result, tbl)
	}
	return result, true
}
