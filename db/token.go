package db

import (
	"fmt"
	dblayer "github.com/zoujiepro/file-server/db/mysql"
)

func GetUsernameByToken(token string) (string, bool) {
	sqlStr := `select user_name from tbl_user_token where user_token = ? limit 1`

	fmt.Printf("sql : %s\n parameter: token = %s\n", sqlStr, token)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("sql prepare err: ", err.Error())
		return "", false
	}

	defer stmt.Close()

	var userName string
	err = stmt.QueryRow(token).Scan(&userName)

	if err != nil {
		fmt.Printf("sql query err: ", err.Error())
		return "", false
	}

	return userName, true
}
