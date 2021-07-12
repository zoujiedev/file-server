package handler

import (
	"fmt"
	"github.com/zoujiepro/file-server/utils"
	"net/http"
)

func HTTPInterceptor(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		if "" == token {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var username string
		var success bool
		if username, success = utils.CheckToken(token); !success {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Printf("current login user: %s\n", username)
		handlerFunc(w, r)
	}
}
