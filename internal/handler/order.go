package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/service"
)

type OrderHandler struct {
	orderService service.OrderServiceInterface
}

func NewOrderHandler(orderService service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (o *OrderHandler) GetOrders(res http.ResponseWriter, req *http.Request) {
	user, err := req.Cookie("user")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	orders, error := o.orderService.GetOrders(user.Value)
	if error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if orders == nil {
		//нет данных для ответа
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

func (o *OrderHandler) AddOrders(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	number, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	isValid := helpers.LuhnAlgorithm(string(number))
	if !isValid {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	existingOrder, err := o.orderService.GetOrder(string(number))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := req.Cookie("user")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if existingOrder != nil {
		if existingOrder.Login == user.Value {
			res.WriteHeader(http.StatusOK)
			return
		} else {
			res.WriteHeader(http.StatusConflict)
		}
		return
	}

	//get status and accural form foreign service
	order := model.Order{Login: user.Value, Number: string(number)}
	err = o.orderService.AddOrder(order)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}
