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
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	godotenv.Load(".env.bank_c")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bankID := os.Getenv("BANK_ID")

	queryForUpdateStatus := `UPDATE transactions 
            SET update_at = extract(epoch from now())::bigint, status = 1
            WHERE update_at < extract(epoch from now())::bigint - 5 AND status = 2`
	go generator.UpdateStatus(ctx, pool, generator.UpdateStatusInterval, queryForUpdateStatus)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			t := generator.NewTransaction(bankID)
			query := `INSERT INTO transactions (uuid, bank_id, from_account, to_account, amount, currency, unixtime, update_at, status)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

			timestamp := t.Timestamp.Unix()

			var status models.StatusCode

			switch {
			case t.Status == models.StatusSuccess:
				status = models.StatusCodeSuccess
			case t.Status == models.StatusFailed:
				status = models.StatusCodeFailed
			case t.Status == models.StatusPending:
				status = models.StatusCodePending
			default:
				log.Fatalf("unknown status: %v", t.Status)
			}

			var currency models.CurrencyISO

			switch {
			case t.Currency == models.RUB:
				currency = models.RUBISO
			case t.Currency == models.EUR:
				currency = models.EURISO
			case t.Currency == models.USD:
				currency = models.USDISO
			default:
				log.Fatalf("unknown currency: %v", t.Currency)
			}

			tag, err := pool.Exec(ctx, query, t.UUID, t.BankID, t.SenderAccount, t.ReceiverAccount, t.Amount, currency, timestamp, timestamp, status)
			time.Sleep(time.Second)
			if err != nil {
				log.Printf("exec error: %v", err)
			} else {
				fmt.Printf("Успешно вставлено строк: %d\n", tag.RowsAffected())
			}
		}
	}
}
