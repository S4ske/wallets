package postgres

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"wallets/internal/domain"
	"wallets/internal/repositories"
)

type WalletRepository struct {
	db *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (wr *WalletRepository) GetWallet(id string) (*domain.Wallet, error) {
	query := "SELECT id, balance FROM wallets WHERE id = $1"
	row := wr.db.QueryRow(query, id)
	var wallet domain.Wallet
	if err := row.Scan(&wallet.ID, &wallet.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrWalletNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

func (wr *WalletRepository) GetWallets() ([]*domain.Wallet, error) {
	query := "SELECT id, balance FROM wallets"
	rows, err := wr.db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	wallets := make([]*domain.Wallet, 0)
	for rows.Next() {
		var w domain.Wallet
		if err := rows.Scan(&w.ID, &w.Balance); err != nil {
			return nil, err
		}
		wallets = append(wallets, &w)
	}
	return wallets, nil
}

func (wr *WalletRepository) SaveWallet(balance int) error {
	query := "INSERT INTO wallets (balance) VALUES ($1)"
	_, err := wr.db.Exec(query, balance)
	if err != nil {
		var pqErr pq.Error
		if errors.As(err, &pqErr) && (pqErr.Code == "23514" || pqErr.Code == "23000") {
			return repositories.ErrBalanceConstraint
		}
		return err
	}
	return err
}

func (wr *WalletRepository) UpdateWallet(ID string, balanceDelta int) error {
	tx, err := wr.db.Begin()
	if err != nil {
		return err
	}
	blockQuery := "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE"
	_, err = tx.Exec(blockQuery, ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	query := "UPDATE wallets SET balance = balance + $2 WHERE id = $1"
	_, err = tx.Exec(query, ID, balanceDelta)
	if err != nil {
		tx.Rollback()
		var pqErr pq.Error
		if errors.As(err, &pqErr) && (pqErr.Code == "23514" || pqErr.Code == "23000") {
			return repositories.ErrBalanceConstraint
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
