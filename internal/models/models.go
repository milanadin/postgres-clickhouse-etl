package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type Currency string

const (
	RUB Currency = "RUB"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

type CurrencyISO int

const (
	RUBISO CurrencyISO = 643
	USDISO CurrencyISO = 840
	EURISO CurrencyISO = 978
)

type Status string

const (
	StatusPending Status = "Pending"
	StatusSuccess Status = "Success"
	StatusFailed  Status = "Failed"
)

type StatusExternal string

const (
	StatusCompleted    StatusExternal = "Completed"
	StatusInProcessing StatusExternal = "In processing"
	StatusRejected     StatusExternal = "Rejected"
)

type StatusCode int

const (
	StatusCodePending StatusCode = 2
	StatusCodeSuccess StatusCode = 1
	StatusCodeFailed  StatusCode = 0
)

type Transaction struct {
	UUID            uuid.UUID
	BankID          string
	SenderAccount   string
	ReceiverAccount string
	Amount          decimal.Decimal
	Currency        Currency
	Timestamp       time.Time
	UpdateAt        time.Time
	Status          Status
}

type TransactionBankA struct {
	UUID            uuid.UUID
	BankID          string
	SenderAccount   string
	ReceiverAccount string
	Amount          decimal.Decimal
	Currency        string
	OccurredAt      time.Time
	UpdateAt        time.Time
	Status          string
}

type TransactionBankB struct {
	UUID            uuid.UUID
	BankID          string
	FromAccount     string
	ToAccount       string
	Amount          int64
	Currency        string
	TransactionTime time.Time
	UpdateAt        time.Time
	Status          string
}

type TransactionBankC struct {
	UUID        uuid.UUID
	BankID      string
	FromAccount string
	ToAccount   string
	Amount      decimal.Decimal
	Currency    int
	Unixtime    int64
	UpdateAt    int64
	Status      int
}
