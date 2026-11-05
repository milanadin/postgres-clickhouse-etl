package extractor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ExtractorBankA(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, cursorUUID uuid.UUID, limit int, interval time.Duration) <-chan models.TransactionBankA {
	out := make(chan models.TransactionBankA)
	go func() {
		defer close(out)
		var rows []models.TransactionBankA
		var err error
		for {
			rows, cursor, cursorUUID, err = ReadBankA(ctx, pool, cursor, cursorUUID, limit)

			if err != nil {
				log.Printf("bank A reading error %v", err)
			}
			for _, t := range rows {
				select {
				case <-ctx.Done():
					return
				case out <- t:
				}
			}
			if ctx.Err() != nil {
				return
			}
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}

		}
	}()
	return out
}

func ExtractorBankB(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, cursorUUID uuid.UUID, limit int, interval time.Duration) <-chan models.TransactionBankB {
	out := make(chan models.TransactionBankB)
	go func() {
		defer close(out)
		var rows []models.TransactionBankB
		var err error
		for {
			rows, cursor, cursorUUID, err = ReadBankB(ctx, pool, cursor, cursorUUID, limit)

			if err != nil {
				log.Printf("bank B reading error %v", err)
			}
			for _, t := range rows {
				select {
				case <-ctx.Done():
					return
				case out <- t:
				}

			}
			if ctx.Err() != nil {
				return
			}
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}

		}
	}()
	return out
}

func ExtractorBankC(ctx context.Context, pool *pgxpool.Pool, cursor int64, cursorUUID uuid.UUID, limit int, interval time.Duration) <-chan models.TransactionBankC {
	out := make(chan models.TransactionBankC)
	go func() {
		defer close(out)
		var rows []models.TransactionBankC
		var err error
		for {
			rows, cursor, cursorUUID, err = ReadBankC(ctx, pool, cursor, cursorUUID, limit)

			if err != nil {
				log.Printf("bank C reading error %v", err)
			}
			for _, t := range rows {
				out <- t
			}
			if ctx.Err() != nil {
				return
			}
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}

		}
	}()
	return out
}
