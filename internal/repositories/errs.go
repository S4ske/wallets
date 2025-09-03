package repositories

import "errors"

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrBalanceConstraint = errors.New("balance must be gte 0")
)
