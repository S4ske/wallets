package service

import (
	"errors"
	"wallets/internal/domain"
	"wallets/internal/repositories"
)

var (
	ErrInvalidID           = errors.New("invalid id")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be gte 0")
	ErrNegativeBalance     = errors.New("balance must be gte 0")
)

type WalletSaver interface {
	SaveWallet(balance int) error
}

type WalletGetter interface {
	GetWallet(id string) (*domain.Wallet, error)
	GetWallets() ([]*domain.Wallet, error)
}

type WalletUpdater interface {
	UpdateWallet(ID string, balanceDelta int) error
}

type WalletService struct {
	walletSaver   WalletSaver
	walletGetter  WalletGetter
	walletUpdater WalletUpdater
}

func NewWalletService(
	walletSaver WalletSaver,
	walletGetter WalletGetter,
	walletUpdater WalletUpdater,
) *WalletService {
	return &WalletService{
		walletGetter:  walletGetter,
		walletUpdater: walletUpdater,
		walletSaver:   walletSaver,
	}
}

func (ws *WalletService) CreateNewWallet(balance int) error {
	err := ws.walletSaver.SaveWallet(balance)
	if err != nil {
		if errors.Is(err, repositories.ErrBalanceConstraint) {
			return ErrNegativeBalance
		}
		return err
	}
	return nil
}

func (ws *WalletService) Deposit(walletID string, amount int) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	return ws.updateWallet(walletID, amount)
}

func (ws *WalletService) Withdraw(walletID string, amount int) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	return ws.updateWallet(walletID, -amount)
}

func (ws *WalletService) updateWallet(walletID string, balanceDelta int) error {
	wallet, err := ws.walletGetter.GetWallet(walletID)
	if err != nil {
		if errors.Is(err, repositories.ErrWalletNotFound) {
			return ErrInvalidID
		}
		return err
	}
	if balanceDelta < 0 && wallet.Balance+balanceDelta < 0 {
		return ErrInsufficientBalance
	}
	if err := ws.walletUpdater.UpdateWallet(walletID, balanceDelta); err != nil {
		if errors.Is(err, repositories.ErrBalanceConstraint) {
			return ErrInsufficientBalance
		}
		return err
	}
	return nil
}

func (ws *WalletService) GetWallet(walletID string) (*domain.Wallet, error) {
	wallet, err := ws.walletGetter.GetWallet(walletID)
	if err != nil {
		if errors.Is(err, repositories.ErrWalletNotFound) {
			return nil, ErrInvalidID
		}
		return nil, err
	}
	return wallet, nil
}

func (ws *WalletService) GetWallets() ([]*domain.Wallet, error) {
	wallet, err := ws.walletGetter.GetWallets()
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
