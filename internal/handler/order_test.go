package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fickleDude/gophemart/internal/handler"
	"github.com/fickleDude/gophemart/internal/mocks"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderHandler_GetOrders(t *testing.T) {
	//init handler
	service := mocks.MockOrderServiceInterface{}
	service.On("GetOrders", "1").Return(nil, fmt.Errorf("test"))
	uploadedAt := time.Date(2026, time.May, 25, 13, 51, 12, 6293, time.UTC)
	service.On("GetOrders", "2").Return([]*model.Order{
		{Login: "2", Number: "1", Status: "PROCESSED", Accrual: 500, UploadedAt: uploadedAt.Format(time.RFC3339)},
	}, nil)

	service.On("GetOrders", "3").Return(nil, nil)

	handler := handler.NewOrderHandler(&service)
	//test table
	type want struct {
		code     int
		request  string
		response string
	}
	tests := []struct {
		name string
		user string
		want want
	}{
		{
			name: "status internal error",
			user: "1",
			want: want{
				code:     500,
				request:  "/api/user/orders",
				response: ``,
			},
		},
		{
			name: "status ok",
			user: "2",
			want: want{
				code:    200,
				request: "/api/user/orders",
				response: `[
								{
									"number": "1",
									"status": "PROCESSED",
									"accrual": 500,
									"uploaded_at": "2026-05-25T13:51:12Z"
								}
							]`,
			},
		},
		{
			name: "status no content",
			user: "3",
			want: want{
				code:     204,
				request:  "/api/user/orders",
				response: ``,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.want.request, nil)
			request.Header.Set("Authorization", test.user)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetOrders)
			h(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode == 200 {
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, res.Header.Get("Content-Type"), "application/json")
				assert.JSONEq(t, test.want.response, string(resBody))
			}
		})
	}
}

func TestOrderHandler_AddOrders(t *testing.T) {
	//init handler
	service := mocks.MockOrderServiceInterface{}
	//неверный формат номера заказа
	service.On("ValidateOrder", "aaa").Return(fmt.Errorf("test"))
	//номер заказа уже был загружен другим пользователем
	service.On("ValidateOrder", "1").Return(nil)
	service.On("GetOrder", "1").Return(&model.Order{Login: "1", Number: "1"}, nil)
	//номер заказа уже был загружен этим пользователем
	service.On("ValidateOrder", "2").Return(nil)
	service.On("GetOrder", "2").Return(&model.Order{Login: "4", Number: "2"}, nil)
	//новый номер заказа принят в обработку
	service.On("ValidateOrder", "3").Return(nil)
	service.On("GetOrder", "3").Return(nil, nil)
	service.On("AddOrder", model.Order{Login: "2", Number: "3"}).Return(nil)
	//внутренняя ошибка сервера
	service.On("ValidateOrder", "4").Return(nil)
	service.On("GetOrder", "4").Return(nil, nil)
	service.On("AddOrder", model.Order{Login: "5", Number: "4"}).Return(fmt.Errorf("test"))

	handler := handler.NewOrderHandler(&service)
	//test table
	type want struct {
		code               int
		request            string
		requestBody        string
		requestContentType string
	}
	tests := []struct {
		name string
		user string
		want want
	}{
		//неверный формат запроса
		{
			name: "status bad request",
			user: "1",
			want: want{
				code:               400,
				request:            "/api/user/orders",
				requestContentType: "application/json",
			},
		},
		//неверный формат номера заказа
		{
			name: "status unprocessable entity",
			user: "1",
			want: want{
				code:               422,
				request:            "/api/user/orders",
				requestBody:        `aaa`,
				requestContentType: "text/plain",
			},
		},
		//номер заказа уже был загружен другим пользователем
		{
			name: "status conflict",
			user: "3",
			want: want{
				code:               409,
				request:            "/api/user/orders",
				requestBody:        `1`,
				requestContentType: "text/plain",
			},
		},
		//номер заказа уже был загружен этим пользователем
		{
			name: "status ok",
			user: "4",
			want: want{
				code:               200,
				request:            "/api/user/orders",
				requestBody:        `2`,
				requestContentType: "text/plain",
			},
		},
		//новый номер заказа принят в обработку
		{
			name: "status accepted",
			user: "2",
			want: want{
				code:               202,
				request:            "/api/user/orders",
				requestBody:        `3`,
				requestContentType: "text/plain",
			},
		},
		//внутренняя ошибка сервера
		{
			name: "status internal server error",
			user: "5",
			want: want{
				code:               500,
				request:            "/api/user/orders",
				requestBody:        `4`,
				requestContentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.request, strings.NewReader(test.want.requestBody))
			request.Header.Set("Authorization", test.user)
			request.Header.Set("Content-Type", test.want.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.AddOrders)
			h(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
