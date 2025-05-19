package postgress

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/wallet/internal/models"
)

type WalletRepository struct {
	DB *sqlx.DB
}

func (wr *WalletRepository) UpdateBalance(id uuid.UUID, amount int, opType models.OperationType) error {
	var query string
	//fmt.Println(opType)
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	tx, err := wr.DB.Begin()
	if err != nil {
		return fmt.Errorf("fail transaction: %w", err)
	}
	fmt.Println("trans")

	defer tx.Rollback()

	fmt.Println("check ok")
	switch opType {
	case "DEPOSIT":
		query = "UPDATE wallets SET balance = balance + $1 WHERE id = $2 RETURNING balance"
	case "WITHDRAW":
		query = "UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1 RETURNING balance"
	default:
		return fmt.Errorf("invalid op type")
	}
	fmt.Println("case")

	var newBalance int
	err = tx.QueryRow(query, amount, id).Scan(&newBalance)
	if err != nil {
		if err == sql.ErrNoRows && opType == "WITHDRAW" {
			return fmt.Errorf("insufficient funds")
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}
	fmt.Println("get balance", newBalance)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	fmt.Println("commit")
	return err
}

func (wr *WalletRepository) GetBalance(id string) (int64, error) {
	var balance int64
	err := wr.DB.QueryRow("SELECT balance FROM wallets WHERE id = $1", id).Scan(&balance)
	return balance, err
}

// Вызываем метод функ. запуска транзакции и если ловим ошибку, пробуем снова
func (wr *WalletRepository) UpdateBalanceWithRetry(id uuid.UUID, amount int, opType models.OperationType, maxRetries int) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := wr.UpdateBalance(id, amount, opType)
		if err == nil {
			return nil
		}

		// Если словили дедлок
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "40P01" {
			lastErr = err
			delay := time.Duration(math.Pow(2, float64(i))) * time.Millisecond
			delay += time.Duration(rand.Intn(100)) * time.Millisecond // Добавляем случайность

			time.Sleep(delay)
			continue
		}

		return err
	}

	return fmt.Errorf("max retries (%d) reached, last error: %v", maxRetries, lastErr)
}
