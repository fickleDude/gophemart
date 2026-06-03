package mocks

import (
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockBalanceService struct {
	mock.Mock
}

func (m *MockBalanceService) GetBalance(login string) (*model.Balance, error) {
	args := m.Called(login)
	return args.Get(0).(*model.Balance), args.Error(1)
}
