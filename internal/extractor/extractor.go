package extractor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ExtractorBankA(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, limit int, interval time.Duration) <-chan TransactionBankA {
	out := make(chan TransactionBankA)
	go func() {
		defer close(out)
		var rows []TransactionBankA
		var err error
		for {
			rows, cursor, err = ReadBankA(ctx, pool, cursor, limit)

			if err != nil {
				log.Printf("bank A reading error %v", err)
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

func ExtractorBankB(ctx context.Context, pool *pgxpool.Pool, cursor time.Time, limit int, interval time.Duration) <-chan TransactionBankB {
	out := make(chan TransactionBankB)
	go func() {
		defer close(out)
		var rows []TransactionBankB
		var err error
		for {
			rows, cursor, err = ReadBankB(ctx, pool, cursor, limit)

			if err != nil {
				log.Printf("bank B reading error %v", err)
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

func ExtractorBankC(ctx context.Context, pool *pgxpool.Pool, cursor int64, cursorUUID uuid.UUID, limit int, interval time.Duration) <-chan TransactionBankC {
	out := make(chan TransactionBankC)
	go func() {
		defer close(out)
		var rows []TransactionBankC
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
