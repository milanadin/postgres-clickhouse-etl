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
)

func main() {
	ctx := context.Background()
	godotenv.Load(".env.bank_c")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bankID := os.Getenv("BANK_ID")

	for {
		t := generator.NewTransaction(bankID)
		query := `INSERT INTO transactions (uuid, bank_id, from_account, to_account, amount, currency, unixtime, status)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		timestamp := t.Timestamp.Unix()

		var status generator.StatusCode

		switch {
		case t.Status == generator.StatusSuccess:
			status = generator.StatusCodeSuccess
		case t.Status == generator.StatusFailed:
			status = generator.StatusCodeFailed
		case t.Status == generator.StatusPending:
			status = generator.StatusCodePending
		}

		var currency generator.CurrencyISO

		switch {
		case t.Currency == generator.RUB:
			currency = generator.RUBISO
		case t.Currency == generator.EUR:
			currency = generator.EURISO
		case t.Currency == generator.USD:
			currency = generator.USDISO
		}

		tag, err := pool.Exec(ctx, query, t.UUID, t.BankID, t.SenderAccount, t.ReceiverAccount, t.Amount, currency, timestamp, status)
		if err != nil {
			log.Printf("exec error: %v", err)
		}
		fmt.Printf("Успешно вставлено строк: %d\n", tag.RowsAffected())
		time.Sleep(time.Second)
	}
}
