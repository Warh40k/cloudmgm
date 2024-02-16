package middleware

import (
	"github.com/Warh40k/cloud-manager/internal/api/service/utils"
	"net/http"
	"strings"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		authHeader := r.Header.Get("Authorization")
		headSplit := strings.Split(authHeader, "Bearer ")
		if len(headSplit) == 2 {
			token = headSplit[1]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := utils.CheckJWT(token)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
