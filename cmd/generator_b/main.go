package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/generator"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
	"github.com/shopspring/decimal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	godotenv.Load(".env.bank_b")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bankID := os.Getenv("BANK_ID")

	queryForUpdateStatus := `UPDATE transactions 
            SET update_at = now(), status = 'Completed'
            WHERE update_at < now() - interval '5 seconds' AND status = 'In processing'`
	go generator.UpdateStatus(ctx, pool, generator.UpdateStatusInterval, queryForUpdateStatus)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			t := generator.NewTransaction(bankID)
			query := `INSERT INTO transactions (uuid, bank_id, from_account, to_account, amount, currency, transaction_time, update_at, status)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			amountInCents := t.Amount.Mul(decimal.NewFromInt(100)).IntPart()

			var status models.StatusExternal

			switch {
			case t.Status == models.StatusSuccess:
				status = models.StatusCompleted
			case t.Status == models.StatusPending:
				status = models.StatusInProcessing
			case t.Status == models.StatusFailed:
				status = models.StatusRejected
			default:
				log.Fatalf("unknown status: %v", t.Status)
			}
			tag, err := pool.Exec(ctx, query, t.UUID, t.BankID, t.SenderAccount, t.ReceiverAccount, amountInCents, t.Currency, t.Timestamp, t.Timestamp, status)
			time.Sleep(time.Second)
			if err != nil {
				log.Printf("exec error: %v", err)
			} else {
				fmt.Printf("Успешно вставлено строк: %d\n", tag.RowsAffected())
			}
		}
	}
}
