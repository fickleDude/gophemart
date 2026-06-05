package middleware

import (
	"net/http"
	"strings"

	"github.com/fickleDude/gophemart/internal/helpers"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuth := strings.Contains(r.URL.Path, "register") || strings.Contains(r.URL.Path, "login")
		if !isAuth {
			cookie, err := helpers.GetCookie(r, "token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			err = helpers.ValidateJWTToken(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
