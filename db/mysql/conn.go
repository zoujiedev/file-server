package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "xxx:xxx@tcp(ip:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("fail to connect to mysql, err: %s", err.Error())
		os.Exit(1)
	}
}

func DbConnect() *sql.DB {
	return db
}
