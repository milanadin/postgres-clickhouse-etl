package extractor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shopspring/decimal"
)

type TransactionBankA struct {
	UUID            uuid.UUID
	BankID          string
	SenderAccount   string
	ReceiverAccount string
	Amount          decimal.Decimal
	Currency        string
	OccurredAt      time.Time
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
	Status      int
}

func ReadBankA(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, limit int) ([]TransactionBankA, time.Time, error) {
	out := []TransactionBankA{}

	last := cursor

	query := `SELECT uuid, bank_id, sender_account, receiver_account, amount, currency, occurred_at, status
			  FROM transactions
	 		  WHERE occurred_at > $1 
			  ORDER BY occurred_at 
			  LIMIT $2`
	rows, err := pool.Query(ctx, query, cursor, limit)
	if err != nil {
		return nil, cursor, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TransactionBankA
		if err := rows.Scan(&t.UUID, &t.BankID, &t.SenderAccount, &t.ReceiverAccount, &t.Amount, &t.Currency, &t.OccurredAt, &t.Status); err != nil {
			return nil, cursor, err
		}
		out = append(out, t)
		last = t.OccurredAt
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, err
	}

	return out, last, nil
}

func ReadBankB(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, limit int) ([]TransactionBankB, time.Time, error) {
	out := []TransactionBankB{}

	last := cursor

	query := `SELECT uuid, bank_id, from_account, to_account, amount, currency, transaction_time, status
			  FROM transactions
	 		  WHERE transaction_time > $1 
			  ORDER BY transaction_time 
			  LIMIT $2`
	rows, err := pool.Query(ctx, query, cursor, limit)
	if err != nil {
		return nil, cursor, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TransactionBankB
		if err := rows.Scan(&t.UUID, &t.BankID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Currency, &t.TransactionTime, &t.Status); err != nil {
			return nil, cursor, err
		}
		out = append(out, t)
		last = t.TransactionTime
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, err
	}

	return out, last, nil
}

func ReadBankC(ctx context.Context, pool *pgxpool.Pool, cursor int64, cursorUUID uuid.UUID, limit int) ([]TransactionBankC, int64, uuid.UUID, error) {
	out := []TransactionBankC{}

	last := cursor
	lastUUID := cursorUUID

	query := `SELECT uuid, bank_id, from_account, to_account, amount, currency, unixtime, status
			  FROM transactions
	 		  WHERE (unixtime, uuid) > ($1, $2)
			  ORDER BY unixtime, uuid
			  LIMIT $3`
	rows, err := pool.Query(ctx, query, cursor, cursorUUID, limit)
	if err != nil {
		return nil, cursor, cursorUUID, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TransactionBankC
		if err := rows.Scan(&t.UUID, &t.BankID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Currency, &t.Unixtime, &t.Status); err != nil {
			return nil, cursor, cursorUUID, err
		}
		out = append(out, t)
		last = t.Unixtime
		lastUUID = t.UUID
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, cursorUUID, err
	}

	return out, last, lastUUID, nil
}
