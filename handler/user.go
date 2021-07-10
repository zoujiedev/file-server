package handler

import (
	"fmt"
	"github.com/zoujiepro/file-server/db"
	"github.com/zoujiepro/file-server/utils"
	"net/http"
	"time"
)

const (
	passwordSalt = "#880!"
	tokenSalt    = "_token_salt_@"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		utils.WriteFail(w, "invalid parameters")
		return
	}

	encPassword := utils.Sha1Bytes([]byte(password + passwordSalt))

	if success := db.UserSignup(username, encPassword); success {
		utils.WriteSuccess(w, nil)
	} else {
		utils.WriteFail(w, "sign up fail")
	}
}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	//校验参数合法性
	if len(username) < 3 || len(password) < 5 {
		utils.WriteFail(w, "username or password is not invalid")
		return
	}

	//检验用户合法性
	if success := db.UserSignIn(username, utils.Sha1Bytes([]byte(password+passwordSalt))); !success {
		utils.WriteFail(w, "sign in fail")
		return
	}

	//生成token
	token := getToken(username)
	db.UpdateUserToken(username, token)

	utils.WriteSuccess(w, token)
}

func getToken(username string) string {
	//40字符 MD5(username+timestamp+token_salt)+[:8]timestamp
	ts := fmt.Sprintf("%x", time.Now().Unix())

	tokenPrefix := utils.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}
