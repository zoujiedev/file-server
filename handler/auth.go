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

		if username, success := utils.CheckToken(token); success {
			fmt.Printf("current login user: %s\n", username)
			handlerFunc(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}
}
