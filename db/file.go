package db

import (
	"database/sql"
	"fmt"
	mydb "github.com/zoujiepro/file-server/db/mysql"
)

func UploadFile(fileHash string, fileName string, fileSize int64, fileAddr string) bool {
	stmt, err := mydb.DbConnect().Prepare(`insert into tbl_file(file_sha1,file_name,file_size,file_addr,status) values (?,?,?,?,?)`)

	if err != nil {
		fmt.Printf("fail to prepare statment, err: %s", err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr, 1)
	if err != nil {
		fmt.Printf("fail to exec sql, err: %s", err.Error())
		return false
	}

	if affect, err := ret.RowsAffected(); err == nil && affect > 0 {
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DbConnect().Prepare(`select file_sha1,file_name,file_size,file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1`)
	if err != nil {
		fmt.Printf("prepare sql err: %s", err.Error())
		return nil, err
	}

	defer stmt.Close()

	tbf := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tbf.FileHash, &tbf.FileName, &tbf.FileSize, &tbf.FileAddr)
	if err != nil {
		fmt.Printf("query sql err: %s", err.Error())
		return nil, err
	}

	return &tbf, nil
}
