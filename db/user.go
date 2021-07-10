package db

import (
	"database/sql"
	"fmt"
	dblayer "github.com/zoujiepro/file-server/db/mysql"
)

func UserSignup(userName string, userPassword string) bool {
	sqlStr := `insert into tbl_user(user_name,user_pwd) values (?,?)`

	fmt.Printf("the sql: %s\n the parameter: userName = %s, userPassword = %s\n", sqlStr, userName, userPassword)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
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

func UserSignIn(username string, encUserPassword string) bool {

	sqlStr := `select user_pwd from tbl_user where user_name = ? limit 1`

	fmt.Printf("the sql: %s\n the parameter: userName = %s, encUserPassword = %s\n", sqlStr, username, encUserPassword)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("fail to Prepare, err: %s", err.Error())
		return false
	}

	defer stmt.Close()

	row := stmt.QueryRow(username)
	if err != nil {
		fmt.Println(err)
		return false
	} else if row == nil {
		fmt.Println("not find user: " + username)
		return false
	}

	var sqlPassword sql.NullString
	err = row.Scan(&sqlPassword)
	if err != nil {
		fmt.Println("row scan err: " + err.Error())
		return false
	}

	return len(sqlPassword.String) > 0 && sqlPassword.String == encUserPassword
}

func UpdateUserToken(username string, userToken string) bool {
	sqlStr := `replace into tbl_user_token (user_name,user_token) values (?,?)`

	fmt.Printf("the sql: %s\n the parameter: userName = %s, userToken = %s\n", sqlStr, username, userToken)

	stmt, err := dblayer.DbConnect().Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare sql err: " + err.Error())
		return false
	}

	result, err := stmt.Exec(username, userToken)
	if err != nil {
		fmt.Println("update err: " + err.Error())
		return false
	}

	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return true
	}
	return false
}
