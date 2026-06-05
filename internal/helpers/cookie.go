package helpers

import "net/http"

func SetCookie(res http.ResponseWriter, name string, value string) {
	token := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,                    // Доступ только через HTTP, защита от XSS
		Secure:   true,                    // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(res, token)
}

func GetCookie(req *http.Request, name string) (*http.Cookie, error) {
	return req.Cookie(name)
}
