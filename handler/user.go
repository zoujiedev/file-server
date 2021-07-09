package handler

import (
	"github.com/zoujiepro/file-server/db"
	"github.com/zoujiepro/file-server/utils"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	passwordSalt = "#880!"
)

func UserSign(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		file, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(file)
		return
	}

	if "POST" == r.Method {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if len(username) < 3 || len(password) < 5 {
			io.WriteString(w, "invalid parameters")
			return
		}

		encPassword := utils.Sha1Bytes([]byte(password + passwordSalt))

		if success := db.UserSign(username, encPassword); success {
			io.WriteString(w, "success")
		} else {
			io.WriteString(w, "fail")
		}

		return
	}
}
