package service_tests

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"wallets/internal/domain"
	"wallets/internal/repositories"
	"wallets/internal/service"
)

type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) SaveWallet(balance int) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *MockWalletRepo) GetWallet(id string) (*domain.Wallet, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletRepo) GetWallets() ([]*domain.Wallet, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Wallet), args.Error(1)
}

func (m *MockWalletRepo) UpdateWallet(ID string, balanceDelta int) error {
	args := m.Called(ID, balanceDelta)
	return args.Error(0)
}

func TestCreateNewWallet_HappyPath(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	balance := 1000

	repo.On("SaveWallet", balance).Return(nil)

	err := svc.CreateNewWallet(balance)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestCreateNewWallet_NegativeBalance(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	balance := -1000

	repo.On("SaveWallet", balance).Return(repositories.ErrBalanceConstraint)

	err := svc.CreateNewWallet(balance)

	require.ErrorIs(t, err, service.ErrNegativeBalance)

	repo.AssertExpectations(t)
}

func TestDeposit_HappyPath(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := 1000

	repo.On("GetWallet", id).Return(&domain.Wallet{ID: id, Balance: 0}, nil)
	repo.On("UpdateWallet", id, amount).Return(nil)

	err := svc.Deposit(id, amount)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestDeposit_InvalidID(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "invalidID"
	amount := 1000

	repo.On("GetWallet", id).Return((*domain.Wallet)(nil), repositories.ErrWalletNotFound)

	err := svc.Deposit(id, amount)

	require.ErrorIs(t, err, service.ErrInvalidID)

	repo.AssertExpectations(t)
}

func TestDeposit_NegativeAmount(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := -1000

	err := svc.Deposit(id, amount)

	require.ErrorIs(t, err, service.ErrInvalidAmount)

	repo.AssertExpectations(t)
}

func TestWithdraw_HappyPath(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := 1000

	repo.On("GetWallet", id).Return(&domain.Wallet{ID: id, Balance: 1000}, nil)
	repo.On("UpdateWallet", id, -amount).Return(nil)

	err := svc.Withdraw(id, amount)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestWithdraw_InvalidID(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := 1000

	repo.On("GetWallet", id).Return((*domain.Wallet)(nil), repositories.ErrWalletNotFound)

	err := svc.Withdraw(id, amount)

	require.ErrorIs(t, err, service.ErrInvalidID)

	repo.AssertExpectations(t)
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := 1000

	repo.On("GetWallet", id).Return(&domain.Wallet{ID: id, Balance: 500}, nil)

	err := svc.Withdraw(id, amount)

	require.ErrorIs(t, err, service.ErrInsufficientBalance)

	repo.AssertExpectations(t)
}

func TestWithdraw_NegativeAmount(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"
	amount := -1000

	err := svc.Withdraw(id, amount)

	require.ErrorIs(t, err, service.ErrInvalidAmount)

	repo.AssertExpectations(t)
}

func TestGetWallet_HappyPath(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "neededID"

	repo.On("GetWallet", id).Return(&domain.Wallet{ID: id, Balance: 1000}, nil)

	wallet, err := svc.GetWallet(id)

	require.NoError(t, err)
	require.Equal(t, &domain.Wallet{ID: id, Balance: 1000}, wallet)

	repo.AssertExpectations(t)
}

func TestGetWallet_InvalidID(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	id := "invalidID"

	repo.On("GetWallet", id).Return((*domain.Wallet)(nil), repositories.ErrWalletNotFound)

	wallet, err := svc.GetWallet(id)

	require.ErrorIs(t, err, service.ErrInvalidID)
	require.Equal(t, (*domain.Wallet)(nil), wallet)

	repo.AssertExpectations(t)
}

func TestGetWallets_HappyPath(t *testing.T) {
	repo := new(MockWalletRepo)
	svc := service.NewWalletService(repo, repo, repo)

	repo.On("GetWallets").Return([]*domain.Wallet{{ID: "1", Balance: 1000}, {ID: "2", Balance: 2000}}, nil)

	wallets, err := svc.GetWallets()

	require.NoError(t, err)
	require.Equal(t, []*domain.Wallet{{ID: "1", Balance: 1000}, {ID: "2", Balance: 2000}}, wallets)

	repo.AssertExpectations(t)
}
