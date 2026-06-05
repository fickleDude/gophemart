package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/service"
)

type BalanceHandler struct {
	balanceService service.BalanceServiceInterface
}

func NewBalanceHandler(balanceService service.BalanceServiceInterface) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

func (b *BalanceHandler) GetBalance(res http.ResponseWriter, req *http.Request) {
	//get login from token
	token, err := helpers.GetCookie(req, "token")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	login := helpers.GetUserLogin(token.Value)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := b.balanceService.GetBalance(login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(balance); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	buf.WriteTo(res)
}
