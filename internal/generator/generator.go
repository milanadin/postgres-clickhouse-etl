package generator

import (
	"math/rand/v2"
	"strconv"
	"strings"
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
	Status          Status
}

func accountGenerator() string {

	var builder strings.Builder

	for i := 0; i < 20; i++ {
		digit := rand.IntN(10)
		digitStr := strconv.Itoa(digit)
		builder.WriteString(digitStr)

	}
	return builder.String()
}

func amountGenerator() decimal.Decimal {

	min := 1000
	max := 100000000
	randAmount := rand.IntN(max-min+1) + min
	amount := decimal.New(int64(randAmount), -2)
	return amount
}

func currencyGenerator() Currency {
	switch rand.IntN(3) {
	case 0:
		return RUB
	case 1:
		return EUR
	case 2:
		return USD
	}
	return RUB
}

func statusGenerator() Status {
	n := rand.IntN(100)
	switch {
	case n < 90:
		return StatusSuccess
	case n < 95:
		return StatusPending
	default:
		return StatusFailed
	}
}

func NewTransaction(bankID string) Transaction {

	t := Transaction{
		UUID:            uuid.New(),
		BankID:          bankID,
		SenderAccount:   accountGenerator(),
		ReceiverAccount: accountGenerator(),
		Amount:          amountGenerator(),
		Currency:        currencyGenerator(),
		Timestamp:       time.Now().UTC(),
		Status:          statusGenerator(),
	}
	return t
}
