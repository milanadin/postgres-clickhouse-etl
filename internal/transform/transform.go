package transform

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/extractor"

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

func convertStatusB(status StatusExternal) (Status, bool) {
	switch status {
	case StatusCompleted:
		return StatusSuccess, true
	case StatusInProcessing:
		return StatusPending, true
	case StatusRejected:
		return StatusFailed, true
	default:
		return "", false
	}
}

func convertStatusC(status StatusCode) (Status, bool) {
	switch status {
	case StatusCodeSuccess:
		return StatusSuccess, true
	case StatusCodePending:
		return StatusPending, true
	case StatusCodeFailed:
		return StatusFailed, true
	default:
		return "", false
	}
}

func convertCurrencyC(currency CurrencyISO) (Currency, bool) {
	switch currency {
	case RUBISO:
		return RUB, true
	case USDISO:
		return USD, true
	case EURISO:
		return EUR, true
	default:
		return "", false
	}
}

func convertTimeC(timeUnix int64) time.Time {
	t := time.Unix(timeUnix, 0)
	return t
}

func NormalizeBankA(ctx context.Context, in <-chan extractor.TransactionBankA) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case a, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- Transaction{
					UUID:            a.UUID,
					BankID:          a.BankID,
					SenderAccount:   a.SenderAccount,
					ReceiverAccount: a.ReceiverAccount,
					Amount:          a.Amount,
					Currency:        Currency(a.Currency),
					Timestamp:       a.OccurredAt.UTC(),
					Status:          Status(a.Status),
				}:
				}
			}
		}
	}()
	return out
}

func NormalizeBankB(ctx context.Context, in <-chan extractor.TransactionBankB) <-chan Transaction {
	out := make(chan Transaction)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case b, ok := <-in:
				if !ok {
					return
				}
				status, valid := convertStatusB(StatusExternal(b.Status))
				if !valid {
					log.Printf("bank B: unknown status %q, uuid %s", b.Status, b.UUID)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- Transaction{
					UUID:            b.UUID,
					BankID:          b.BankID,
					SenderAccount:   b.FromAccount,
					ReceiverAccount: b.ToAccount,
					Amount:          decimal.New(b.Amount, -2),
					Currency:        Currency(b.Currency),
					Timestamp:       b.TransactionTime.UTC(),
					Status:          status,
				}:
				}
			}
		}

	}()
	return out
}

func NormalizeBankC(ctx context.Context, in <-chan extractor.TransactionBankC) <-chan Transaction {
	out := make(chan Transaction)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case c, ok := <-in:
				if !ok {
					return
				}
				status, valid := convertStatusC(StatusCode(c.Status))
				if !valid {
					log.Printf("bank C: unknown status %d, uuid %s", c.Status, c.UUID)
					continue
				}
				currency, valid := convertCurrencyC(CurrencyISO(c.Currency))
				if !valid {
					log.Printf("bank C: unknown currency %d, uuid %s", c.Currency, c.UUID)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- Transaction{
					UUID:            c.UUID,
					BankID:          c.BankID,
					SenderAccount:   c.FromAccount,
					ReceiverAccount: c.ToAccount,
					Amount:          c.Amount,
					Currency:        currency,
					Timestamp:       convertTimeC(c.Unixtime).UTC(),
					Status:          status,
				}:
				}
			}
		}

	}()
	return out
}
