package helpers

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var s = securecookie.New(hashKey, blockKey)

func SetCookie(res http.ResponseWriter, name string, value string) {

	if encoded, err := s.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:     name,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(res, cookie)
	}
}

func GetCookie(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err == nil {
		var value string
		if err = s.Decode(name, cookie.Value, &value); err == nil {
			return value, nil
		}
	}
	return "", err
}
