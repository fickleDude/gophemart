package handler

import (
	"net/http"

	"github.com/fickleDude/gophemart/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (o *UserHandler) Login(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("unimplemented"))
}

func (o *UserHandler) Register(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("unimplemented"))
}
