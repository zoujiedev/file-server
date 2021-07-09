package db

import (
	"fmt"
	dblayer "github.com/zoujiepro/file-server/db/mysql"
)

func UserSign(userName string, userPassword string) bool {
	stmt, err := dblayer.DbConnect().Prepare(`insert into tbl_user(user_name,user_pwd) values (?,?)`)
	if err != nil {
		fmt.Printf("fail to Prepare, err: %s", err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(userName, userPassword)
	if err != nil {
		fmt.Printf("fail to insert, err: %s", err.Error())
		return false
	}

	if affect, err := result.RowsAffected(); err == nil && affect > 0 {
		return true
	}

	return false
}
