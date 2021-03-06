package db

import (
	"database/sql"
	"fmt"
	dblayer "github.com/zoujiepro/file-server/db/mysql"
)

func UploadFile(fileHash string, fileName string, fileSize int64, fileAddr string) bool {
	sqlStr := `insert into tbl_file(file_sha1,file_name,file_size,file_addr,status) values (?,?,?,?,?)`

	fmt.Printf("the sql: %s\n the parameter: filehash = %s, fileName = %s, fileSize = %d, fileAddr = %s\n", sqlStr, fileHash, fileName, fileSize, fileAddr)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
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
	sqlStr := `select file_sha1,file_name,file_size,file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1`

	fmt.Printf("the sql: %s\n the parameter: filehash = %s\n", sqlStr, filehash)
	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
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

func DeleteFileMeta(filehash string) bool {
	sqlStr := `delete from tbl_file where file_sha1 = ? limit 1`
	fmt.Printf("the sql: %s\n the parameter: filehash = %s\n", sqlStr, filehash)
	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Println("sql prepare fail: " + err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(filehash)
	if err != nil {
		fmt.Println("sql exec fail: " + err.Error())
		return false
	}

	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return true
	}

	return false
}

func UpdateFileMeta(filehash string, newName string) bool {
	sqlStr := `update tbl_file set file_name = ? where file_sha1 = ?`

	fmt.Printf("the sql: %s\n the parameter: filehash = %s, newName = %s\n", sqlStr, filehash, newName)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("update prepare err: %s", err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(newName, filehash)
	if err != nil {
		fmt.Printf("update exec err: %s", err.Error())
		return false
	}

	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return true
	} else {
		fmt.Printf("sql exec success, but affected %s row", affected)
		return false
	}
}
