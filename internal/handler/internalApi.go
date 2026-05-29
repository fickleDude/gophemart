package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fickleDude/gophemart/internal/service"
	"github.com/go-chi/chi"
)

type InternalApiHandler struct {
	service *service.InternalApiService
}

func NewInternalApiHandler(service *service.InternalApiService) *InternalApiHandler {
	return &InternalApiHandler{service: service}
}

func (i *InternalApiHandler) GetData(res http.ResponseWriter, req *http.Request) {
	number := chi.URLParam(req, "number")
	orders, error := i.service.GetData(number)
	if error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if orders == nil {
		//заказ не зарегистрирован в системе расчёта
		res.WriteHeader(http.StatusNoContent)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(orders); err != nil {
		//внутренняя ошибка сервера
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	buf.WriteTo(res)
}
