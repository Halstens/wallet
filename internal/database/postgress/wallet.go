package postgress

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/wallet/internal/models"
)

type WalletRepository struct {
	DB *sql.DB
}

// func NewWalletRepository(db *sql.DB) *WalletRepository {
// 	return &WalletRepository{db: db}
// }

// func NewStorage() (*Storage, error) {
//     connStr := "user=postgres dbname=wallet password=postgres sslmode=disable"
//     db, err := sql.Open("postgres", connStr)
//     if err != nil {
//         return nil, err
//     }
//     return &Storage{db: db}, nil

func (wr *WalletRepository) UpdateBalance(id uuid.UUID, amount int, opType models.OperationType) error {
	var query string
	fmt.Println(opType)
	switch opType {
	case "DEPOSIT":
		query = "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	case "WITHDRAW":
		query = "UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1"
	default:
		return fmt.Errorf("invalid op type")
	}

	_, err := wr.DB.Exec(query, amount, id)
	return err
}

func (wr *WalletRepository) GetBalance(id string) (int64, error) {

	// stmt := `SELECT balance FROM wallets WHERE id = ?`
	// row := wr.DB.QueryRow(stmt, id)

	// w := &models.Wallet{}

	// err := row.Scan(&w.ID, &w.Balance)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, models.ErrNoRecord
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	var balance int64
	err := wr.DB.QueryRow("SELECT balance FROM wallets WHERE id = $1", id).Scan(&balance)
	return balance, err
	//return w, nil
}

// func (r *WalletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
// 	_, err := r.db.ExecContext(ctx,
// 		"INSERT INTO wallets (id, balance) VALUES ($1, $2)",
// 		wallet.ID, wallet.Balance)
// 	return err
// }

// func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
// 	var wallet model.Wallet
// 	err := r.db.QueryRowContext(ctx,
// 		"SELECT id, balance FROM wallets WHERE id = $1", id).
// 		Scan(&wallet.ID, &wallet.Balance)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &wallet, nil
// }

// func (r *WalletRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) error {
// 	tx, err := r.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	var currentBalance int64
// 	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", id).
// 		Scan(&currentBalance)
// 	if err != nil {
// 		return err
// 	}

// 	newBalance := currentBalance + amount
// 	if newBalance < 0 {
// 		return fmt.Errorf("insufficient funds")
// 	}

// 	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, id)
// 	if err != nil {
// 		return err
// 	}

// 	return tx.Commit()
// }
