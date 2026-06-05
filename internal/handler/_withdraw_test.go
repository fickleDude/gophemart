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

func TestWithdrawHandler_GetWithdraws(t *testing.T) {
	//init handler
	service := mocks.MockWithdrawServiceInterface{}
	service.On("GetWithdraws", "1").Return(nil, fmt.Errorf("test"))
	processedAt := time.Date(2026, time.May, 25, 13, 51, 12, 6293, time.UTC)
	service.On("GetWithdraws", "2").Return([]*model.Withdraw{{
		Order:       "2377225624",
		Sum:         751,
		ProcessedAt: processedAt.Format(time.RFC3339),
	}}, nil)

	service.On("GetWithdraws", "3").Return(nil, nil)

	handler := handler.NewWithdrawHandler(&service, nil)
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
				request:  "/api/user/withdrawals",
				response: ``,
			},
		},
		{
			name: "status ok",
			user: "2",
			want: want{
				code:    200,
				request: "/api/user/withdrawals",
				response: `[
								{
									"order": "2377225624",
									"sum": 751,
									"processed_at": "2026-05-25T13:51:12Z"
								}
							]`,
			},
		},
		{
			name: "status no content",
			user: "3",
			want: want{
				code:     204,
				request:  "/api/user/withdrawals",
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
			h := http.HandlerFunc(handler.GetWithdraws)
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

func TestWithdrawHandler_AddWithdraw(t *testing.T) {
	//init handler
	withdrawService := mocks.MockWithdrawServiceInterface{}
	balanceService := mocks.MockBalanceServiceInterface{}
	withdrawService.On("ValidateOrder", "aaa").Return(fmt.Errorf("test"))
	withdrawService.On("ValidateOrder", "1").Return(nil)
	withdrawService.On("AddWithdraw", model.Withdraw{Login: "2", Order: "1", Sum: 751}).Return(nil)
	balanceService.On("GetBalance", "1").Return(&model.Balance{Current: 0, Withdraw: 0}, nil)
	balanceService.On("GetBalance", "2").Return(&model.Balance{Current: 1000, Withdraw: 0}, nil)

	handler := handler.NewWithdrawHandler(&withdrawService, &balanceService)
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
		//content type
		{
			name: "status bad request",
			user: "1",
			want: want{
				code:               400,
				request:            "/api/user/balance/withdraw",
				requestContentType: "text/html",
			},
		},
		//invalid number
		{
			name: "status unprocessable entity",
			user: "2",
			want: want{
				code:    422,
				request: "/api/user/balance/withdraw",
				requestBody: `{
					"order": "aaa",
					"sum": 751
				}`,
				requestContentType: "application/json",
			},
		},
		//на счету недостаточно средств
		{
			name: "status payment required",
			user: "1",
			want: want{
				code:    402,
				request: "/api/user/balance/withdraw",
				requestBody: `{
					"order": "1",
					"sum": 751
				}`,
				requestContentType: "application/json",
			},
		},
		//invalid json
		{
			name: "status internal server error",
			user: "2",
			want: want{
				code:    500,
				request: "/api/user/balance/withdraw",
				requestBody: `"{
					"order": "1",
					"sum": 751
				}`,
				requestContentType: "application/json",
			},
		},
		//ok
		{
			name: "status ok",
			user: "2",
			want: want{
				code:    200,
				request: "/api/user/balance/withdraw",
				requestBody: `{
					"order": "1",
					"sum": 751
				}`,
				requestContentType: "application/json",
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
			h := http.HandlerFunc(handler.AddWithdraw)
			h(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
