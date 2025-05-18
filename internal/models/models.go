package models

import (
	"errors"

	"github.com/google/uuid"
)

type OperationType string

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

const (
	DEPOSIT  OperationType = "DEPOSIT"
	WITHDRAW OperationType = "WITHDRAW"
)

type Wallet struct {
	ID      uuid.UUID `json:"walletId"`
	Balance int64     `json:"balance"`
}

type WalletOperation struct {
	WalletID      uuid.UUID     `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        int64         `json:"amount"`
}
