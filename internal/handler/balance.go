package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/service"
)

type BalanceHandler struct {
	balanceService service.BalanceServiceInterface
}

func NewBalanceHandler(balanceService service.BalanceServiceInterface) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

func (b *BalanceHandler) GetBalance(res http.ResponseWriter, req *http.Request) {
	user, err := req.Cookie("user")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := b.balanceService.GetBalance(user.Value)
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
