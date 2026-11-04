package loader

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
)

func writeBatch(ctx context.Context, conn clickhouse.Conn, txs []models.Transaction) error {

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO transactions")
	if err != nil {
		return err
	}

	for _, t := range txs {
		err = batch.Append(
			t.UUID,
			t.BankID,
			t.SenderAccount,
			t.ReceiverAccount,
			t.Amount,
			string(t.Currency),
			t.Timestamp,
			t.UpdateAt,
			string(t.Status),
		)
		if err != nil {
			_ = batch.Abort()
			return err
		}
	}
	if err := batch.Send(); err != nil {
		return err
	}

	return nil
}

func Load(ctx context.Context, conn clickhouse.Conn, input <-chan []models.Transaction, maxAttempts int, replay time.Duration) error {
	for {
		var txs []models.Transaction
		select {
		case <-ctx.Done():
			return ctx.Err()
		case batchTx, ok := <-input:
			if !ok {
				return nil
			}
			txs = batchTx
		}

		var err error

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err = writeBatch(ctx, conn, txs)
			if err == nil {
				break
			}

			if ctx.Err() != nil {
				return ctx.Err()
			}

			if attempt == maxAttempts {
				return fmt.Errorf(
					"failed to write batch after %d attempts: %w",
					maxAttempts,
					err,
				)
			}

			timer := time.NewTimer(replay)

			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return ctx.Err()

			case <-timer.C:
			}
		}
	}
}
