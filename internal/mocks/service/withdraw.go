package mocks

import (
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockWithdrawService struct {
	mock.Mock
}

func (m *MockWithdrawService) GetWithdraws(login string) ([]*model.Withdraw, error) {
	args := m.Called(login)
	result, ok := args.Get(0).([]*model.Withdraw)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockWithdrawService) AddWithdraw(withdraw model.Withdraw) error {
	args := m.Called(withdraw)
	return args.Error(0)
}

func (m *MockWithdrawService) ValidateOrder(number string) error {
	args := m.Called(number)
	return args.Error(0)
}
