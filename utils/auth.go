package utils

import "github.com/zoujiepro/file-server/db"

func CheckToken(token string) (string, bool) {
	//todo 有效性校验 过期校验
	return db.GetUsernameByToken(token)
}
