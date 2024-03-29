package middleware

import (
	"context"
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
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		id, err := utils.CheckJWT(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
