package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/generator"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()
	godotenv.Load(".env.bank_b")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bankID := os.Getenv("BANK_ID")

	for {
		t := generator.NewTransaction(bankID)
		query := `INSERT INTO transactions (uuid, bank_id, from_account, to_account, amount, currency, transaction_time, status)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		amountInCents := t.Amount.Mul(decimal.NewFromInt(100)).IntPart()

		var status generator.StatusExternal

		switch {
		case t.Status == generator.StatusSuccess:
			status = generator.StatusCompleted
		case t.Status == generator.StatusPending:
			status = generator.StatusInProcessing
		case t.Status == generator.StatusFailed:
			status = generator.StatusRejected
		}
		tag, err := pool.Exec(ctx, query, t.UUID, t.BankID, t.SenderAccount, t.ReceiverAccount, amountInCents, t.Currency, t.Timestamp, status)
		if err != nil {
			log.Printf("exec error: %v", err)
		}
		fmt.Printf("Успешно вставлено строк: %d\n", tag.RowsAffected())
		time.Sleep(time.Second)
	}
}
