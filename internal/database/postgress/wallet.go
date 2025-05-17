package postgres

import (
	"context"
	"database/sql"
	"fmt"

	model "github.com/wallet/internal/models"

	"github.com/google/uuid"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO wallets (id, balance) VALUES ($1, $2)",
		wallet.ID, wallet.Balance)
	return err
}

func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.QueryRowContext(ctx,
		"SELECT id, balance FROM wallets WHERE id = $1", id).
		Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", id).
		Scan(&currentBalance)
	if err != nil {
		return err
	}

	newBalance := currentBalance + amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
