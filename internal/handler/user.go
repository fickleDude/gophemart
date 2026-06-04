package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Login(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	existingUser, err := u.service.GetUser(user.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.Password != existingUser.Password {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString, err := helpers.CreateJWTToken(user.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,                    // Доступ только через HTTP, защита от XSS
		Secure:   true,                    // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(res, token)
	login := &http.Cookie{
		Name:     "user",
		Value:    user.Login,
		Path:     "/",
		HttpOnly: true,                    // Доступ только через HTTP, защита от XSS
		Secure:   true,                    // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(res, login)
}

func (u *UserHandler) Register(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	existingUser, err := u.service.GetUser(user.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		res.WriteHeader(http.StatusConflict)
		return
	}
	err = u.service.AddUser(user)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/api/user/login", http.StatusTemporaryRedirect)
}
