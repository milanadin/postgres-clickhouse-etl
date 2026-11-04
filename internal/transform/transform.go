package transform

import (
	"context"
	"log"
	"time"

	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"

	"github.com/shopspring/decimal"
)

func convertStatusB(status models.StatusExternal) (models.Status, bool) {
	switch status {
	case models.StatusCompleted:
		return models.StatusSuccess, true
	case models.StatusInProcessing:
		return models.StatusPending, true
	case models.StatusRejected:
		return models.StatusFailed, true
	default:
		return "", false
	}
}

func convertStatusC(status models.StatusCode) (models.Status, bool) {
	switch status {
	case models.StatusCodeSuccess:
		return models.StatusSuccess, true
	case models.StatusCodePending:
		return models.StatusPending, true
	case models.StatusCodeFailed:
		return models.StatusFailed, true
	default:
		return "", false
	}
}

func convertCurrencyC(currency models.CurrencyISO) (models.Currency, bool) {
	switch currency {
	case models.RUBISO:
		return models.RUB, true
	case models.USDISO:
		return models.USD, true
	case models.EURISO:
		return models.EUR, true
	default:
		return "", false
	}
}

func convertTimeC(timeUnix int64) time.Time {
	t := time.Unix(timeUnix, 0)
	return t
}

func NormalizeBankA(ctx context.Context, in <-chan models.TransactionBankA) <-chan models.Transaction {
	out := make(chan models.Transaction)

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
				case out <- models.Transaction{
					UUID:            a.UUID,
					BankID:          a.BankID,
					SenderAccount:   a.SenderAccount,
					ReceiverAccount: a.ReceiverAccount,
					Amount:          a.Amount,
					Currency:        models.Currency(a.Currency),
					Timestamp:       a.OccurredAt.UTC(),
					UpdateAt:        a.UpdateAt.UTC(),
					Status:          models.Status(a.Status),
				}:
				}
			}
		}
	}()
	return out
}

func NormalizeBankB(ctx context.Context, in <-chan models.TransactionBankB) <-chan models.Transaction {
	out := make(chan models.Transaction)
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
				status, valid := convertStatusB(models.StatusExternal(b.Status))
				if !valid {
					log.Printf("bank B: unknown status %q, uuid %s", b.Status, b.UUID)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- models.Transaction{
					UUID:            b.UUID,
					BankID:          b.BankID,
					SenderAccount:   b.FromAccount,
					ReceiverAccount: b.ToAccount,
					Amount:          decimal.New(b.Amount, -2),
					Currency:        models.Currency(b.Currency),
					Timestamp:       b.TransactionTime.UTC(),
					UpdateAt:        b.UpdateAt.UTC(),
					Status:          status,
				}:
				}
			}
		}

	}()
	return out
}

func NormalizeBankC(ctx context.Context, in <-chan models.TransactionBankC) <-chan models.Transaction {
	out := make(chan models.Transaction)
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
				status, valid := convertStatusC(models.StatusCode(c.Status))
				if !valid {
					log.Printf("bank C: unknown status %d, uuid %s", c.Status, c.UUID)
					continue
				}
				currency, valid := convertCurrencyC(models.CurrencyISO(c.Currency))
				if !valid {
					log.Printf("bank C: unknown currency %d, uuid %s", c.Currency, c.UUID)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- models.Transaction{
					UUID:            c.UUID,
					BankID:          c.BankID,
					SenderAccount:   c.FromAccount,
					ReceiverAccount: c.ToAccount,
					Amount:          c.Amount,
					Currency:        currency,
					Timestamp:       convertTimeC(c.Unixtime).UTC(),
					UpdateAt:        convertTimeC(c.UpdateAt).UTC(),
					Status:          status,
				}:
				}
			}
		}

	}()
	return out
}
