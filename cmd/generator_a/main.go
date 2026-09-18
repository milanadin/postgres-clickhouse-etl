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
	godotenv.Load(".env.bank_a")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bankID := os.Getenv("BANK_ID")

	for {
		t := generator.NewTransaction(bankID)
		query := `INSERT INTO transactions (uuid, bank_id, sender_account, receiver_account, amount, currency, occurred_at, status)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		tag, err := pool.Exec(ctx, query, t.UUID, t.BankID, t.SenderAccount, t.ReceiverAccount, t.Amount, t.Currency, t.Timestamp, t.Status)
		if err != nil {
			log.Printf("exec error: %v", err)
		}
		fmt.Printf("Успешно вставлено строк: %d\n", tag.RowsAffected())
		time.Sleep(time.Second)
	}
}
