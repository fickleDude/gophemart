package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/service"
)

type WithdrawHandler struct {
	withdrawService *service.WithdrawService
	balanceService  *service.BalanceService
}

func NewWithdrawHandler(withdrawService *service.WithdrawService, balanceService *service.BalanceService) *WithdrawHandler {
	return &WithdrawHandler{withdrawService: withdrawService, balanceService: balanceService}
}

func (w *WithdrawHandler) GetWithdraws(res http.ResponseWriter, req *http.Request) {
	withdraws, err := w.withdrawService.GetWithdraws("123")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if withdraws == nil {
		res.Write([]byte(" нет ни одного списания"))
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

	err := w.withdrawService.ValidateOrder(withdraw.Order)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	balance, err := w.balanceService.GetBalance("123")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance.Current-withdraw.Sum < 0 {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}
	withdraw.Login = "123"
	err = w.withdrawService.AddWithdraw(withdraw)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
