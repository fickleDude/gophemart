package main

import (
	"net/http"

	"github.com/fickleDude/gophemart/internal/config/db"
	"github.com/fickleDude/gophemart/internal/handler"
	"github.com/fickleDude/gophemart/internal/repository"
	"github.com/fickleDude/gophemart/internal/service"
	"github.com/go-chi/chi"
)

func main() {
	//repository
	storage := db.GetDBConnection()
	defer db.CloseDBConnection()
	orderRepository := repository.NewOrderRepository(storage)
	withdrawRepository := repository.NewWithdrawRepository(storage)
	userRepository := repository.NewUserRepository()
	//services
	orderService := service.NewOrderService(orderRepository)
	withdrawService := service.NewWithdrawService(withdrawRepository)
	balanceService := service.NewBalaneService(orderRepository, withdrawRepository)
	userService := service.NewUserService(userRepository)

	//handlers
	orderHandler := handler.NewOrderHandler(orderService)
	withdrawHandler := handler.NewWithdrawHandler(withdrawService, balanceService)
	userHandler := handler.NewUserHandler(userService)
	balanceHandler := handler.NewBalanceHandler(balanceService)

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			//регистрация пользователя
			r.Post("/register", userHandler.Register)
			//аутентификация пользователя
			r.Post("/login", userHandler.Login)
			r.Route("/orders", func(r chi.Router) {
				//получение списка загруженных пользователем номеров заказов
				r.Get("/", orderHandler.GetOrders)
				//загрузка пользователем номера заказа для расчёта
				r.Post("/", orderHandler.AddOrders)
			})
			r.Route("/withdrawals", func(r chi.Router) {
				//получение информации о выводе средств с накопительного счёта пользователем
				r.Get("/", withdrawHandler.GetWithdraws)
			})
			r.Route("/balance", func(r chi.Router) {
				//получение текущего баланса счёта баллов лояльности пользователя
				r.Get("/", balanceHandler.GetBalance)
				//запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
				r.Post("/withdraw", withdrawHandler.AddWithdraw)
			})
		})
	})

	//start server
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
