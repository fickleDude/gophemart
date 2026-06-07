package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fickleDude/gophemart/internal/handler"
	"github.com/fickleDude/gophemart/internal/mocks"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register(t *testing.T) {
	//init handler
	service := mocks.MockUserServiceInterface{}
	//логин уже занят
	service.On("GetUser", "3").Return(&model.User{Login: "3", Password: "rrty"}, nil)
	//пользователь успешно зарегистрирован и аутентифицирован
	service.On("GetUser", "4").Return(nil, nil)
	service.On("AddUser", model.User{Login: "4", Password: "pass"}).Return(nil, nil)

	handler := handler.NewUserHandler(&service)
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
		// внутренняя ошибка сервера
		{
			name: "status internal error",
			user: "1",
			want: want{
				code:               500,
				request:            "/api/user/register",
				requestContentType: "application/json",
				requestBody:        `{"login : "1", password : "pass"}`,
			},
		},
		//неверный формат запроса
		{
			name: "status bad request",
			user: "2",
			want: want{
				code:               400,
				request:            "/api/user/register",
				requestContentType: "text/plain",
				requestBody:        `{"login : "2", password : "pass"}`,
			},
		},
		//логин уже занят
		{
			name: "status conflict",
			user: "3",
			want: want{
				code:               409,
				request:            "/api/user/register",
				requestContentType: "application/json",
				requestBody:        `{"login" : "3", "password" : "pass"}`,
			},
		},
		//пользователь успешно зарегистрирован и аутентифицирован
		{
			name: "status ok",
			user: "4",
			want: want{
				code:               307, //redirect to login
				request:            "/api/user/register",
				requestContentType: "application/json",
				requestBody:        `{"login" : "4", "password" : "pass"}`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//create request
			request := httptest.NewRequest(http.MethodPost, test.want.request, strings.NewReader(test.want.requestBody))
			request.Header.Set("Content-Type", test.want.requestContentType)
			//add recorder
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Register)
			//execute
			h(w, request)
			//parse
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	//init handler
	service := mocks.MockUserServiceInterface{}
	service.On("ValidateUser", model.User{Login: "3", Password: "pass"}).Return(false, nil)
	service.On("ValidateUser", model.User{Login: "4", Password: "pass"}).Return(true, nil)

	handler := handler.NewUserHandler(&service)
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
		// внутренняя ошибка сервера
		{
			name: "status internal error",
			user: "1",
			want: want{
				code:               500,
				request:            "/api/user/login",
				requestContentType: "application/json",
				requestBody:        `{"login : "1", password : "pass"}`,
			},
		},
		//неверный формат запроса
		{
			name: "status bad request",
			user: "2",
			want: want{
				code:               400,
				request:            "/api/user/login",
				requestContentType: "text/plain",
				requestBody:        `{"login : "2", password : "pass"}`,
			},
		},
		//неверная пара логин/пароль
		{
			name: "status unauthorized",
			user: "3",
			want: want{
				code:               401,
				request:            "/api/user/login",
				requestContentType: "application/json",
				requestBody:        `{"login" : "3", "password" : "pass"}`,
			},
		},
		//пользователь успешно аутентифицирован
		{
			name: "status ok",
			user: "4",
			want: want{
				code:               200,
				request:            "/api/user/login",
				requestContentType: "application/json",
				requestBody:        `{"login" : "4", "password" : "pass"}`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//create request
			request := httptest.NewRequest(http.MethodPost, test.want.request, strings.NewReader(test.want.requestBody))
			request.Header.Set("Content-Type", test.want.requestContentType)
			//add recorder
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Login)
			//execute
			h(w, request)
			//parse
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode == 200 {
				cookies := res.Cookies()
				exists := false
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						exists = true
						return
					}
				}
				assert.True(t, exists)
			}
		})
	}
}
