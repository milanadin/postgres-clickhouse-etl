package extractor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
)

func ReadBankA(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, cursorUUID uuid.UUID, limit int) ([]models.TransactionBankA, time.Time, uuid.UUID, error) {
	out := []models.TransactionBankA{}

	last := cursor
	lastUUID := cursorUUID

	query := `SELECT uuid, bank_id, sender_account, receiver_account, amount, currency, occurred_at, update_at, status
			  FROM transactions
	 		  WHERE (update_at, uuid) > ($1, $2)
			  ORDER BY update_at, uuid
			  LIMIT $3`
	rows, err := pool.Query(ctx, query, cursor, cursorUUID, limit)
	if err != nil {
		return nil, cursor, cursorUUID, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TransactionBankA
		if err := rows.Scan(&t.UUID, &t.BankID, &t.SenderAccount, &t.ReceiverAccount, &t.Amount, &t.Currency, &t.OccurredAt, &t.UpdateAt, &t.Status); err != nil {
			return nil, cursor, cursorUUID, err
		}
		out = append(out, t)
		last = t.UpdateAt
		lastUUID = t.UUID
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, cursorUUID, err
	}

	return out, last, lastUUID, nil
}

func ReadBankB(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, cursorUUID uuid.UUID, limit int) ([]models.TransactionBankB, time.Time, uuid.UUID, error) {
	out := []models.TransactionBankB{}

	last := cursor
	lastUUID := cursorUUID

	query := `SELECT uuid, bank_id, from_account, to_account, amount, currency, transaction_time, update_at, status
			  FROM transactions
	 		  WHERE (update_at, uuid) > ($1, $2)
			  ORDER BY update_at, uuid
			  LIMIT $3`
	rows, err := pool.Query(ctx, query, cursor, cursorUUID, limit)
	if err != nil {
		return nil, cursor, cursorUUID, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TransactionBankB
		if err := rows.Scan(&t.UUID, &t.BankID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Currency, &t.TransactionTime, &t.UpdateAt, &t.Status); err != nil {
			return nil, cursor, cursorUUID, err
		}
		out = append(out, t)
		last = t.UpdateAt
		lastUUID = t.UUID
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, cursorUUID, err
	}

	return out, last, lastUUID, nil
}

func ReadBankC(ctx context.Context, pool *pgxpool.Pool, cursor int64, cursorUUID uuid.UUID, limit int) ([]models.TransactionBankC, int64, uuid.UUID, error) {
	out := []models.TransactionBankC{}

	last := cursor
	lastUUID := cursorUUID

	query := `SELECT uuid, bank_id, from_account, to_account, amount, currency, unixtime, update_at, status
			  FROM transactions
	 		  WHERE (update_at, uuid) > ($1, $2)
			  ORDER BY update_at, uuid
			  LIMIT $3`
	rows, err := pool.Query(ctx, query, cursor, cursorUUID, limit)
	if err != nil {
		return nil, cursor, cursorUUID, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TransactionBankC
		if err := rows.Scan(&t.UUID, &t.BankID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Currency, &t.Unixtime, &t.UpdateAt, &t.Status); err != nil {
			return nil, cursor, cursorUUID, err
		}
		out = append(out, t)
		last = t.UpdateAt
		lastUUID = t.UUID
	}

	if err := rows.Err(); err != nil {
		return nil, cursor, cursorUUID, err
	}

	return out, last, lastUUID, nil
}
