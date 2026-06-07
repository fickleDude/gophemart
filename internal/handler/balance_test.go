package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fickleDude/gophemart/internal/handler"
	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/mocks"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceHandler_GetBalance(t *testing.T) {
	//init handler
	service := mocks.MockBalanceServiceInterface{}
	service.On("GetBalance", "1").Return(nil, fmt.Errorf("test"))
	service.On("GetBalance", "2").Return(&model.Balance{Current: 500.5, Withdraw: 42}, nil)

	handler := handler.NewBalanceHandler(&service)
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
				request:  "/api/user/balance",
				response: ``,
			},
		},
		{
			name: "status ok",
			user: "2",
			want: want{
				code:    200,
				request: "/api/user/balance",
				response: `  {
								"current": 500.5,
								"withdrawn": 42
							}`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.want.request, nil)
			request.Header.Set("Authorization", test.user)
			//login
			token, _ := helpers.CreateJWTToken(test.user)
			helpers.SetRequestCookie(request, "token", token)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetBalance)
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
