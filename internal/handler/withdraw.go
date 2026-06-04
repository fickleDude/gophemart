package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/service"
)

type WithdrawHandler struct {
	withdrawService service.WithdrawServiceInterface
	balanceService  service.BalanceServiceInterface
}

func NewWithdrawHandler(withdrawService service.WithdrawServiceInterface, balanceService service.BalanceServiceInterface) *WithdrawHandler {
	return &WithdrawHandler{withdrawService: withdrawService, balanceService: balanceService}
}

func (w *WithdrawHandler) GetWithdraws(res http.ResponseWriter, req *http.Request) {
	user, err := req.Cookie("user")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	withdraws, err := w.withdrawService.GetWithdraws(user.Value)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if withdraws == nil {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(withdraws); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	buf.WriteTo(res)
}

func (w *WithdrawHandler) AddWithdraw(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var withdraw model.Withdraw
	if err := json.NewDecoder(req.Body).Decode(&withdraw); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	isValid := helpers.LuhnAlgorithm(withdraw.Order)
	if !isValid {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user, err := req.Cookie("user")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := w.balanceService.GetBalance(user.Value)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance.Current-withdraw.Sum < 0 {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	withdraw.Login = user.Value
	err = w.withdrawService.AddWithdraw(withdraw)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
