package generator

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const UpdateStatusInterval = 1 * time.Second

func UpdateStatus(ctx context.Context, pool *pgxpool.Pool, interval time.Duration, query string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := pool.Exec(ctx, query)
			if err != nil {
				log.Printf("unable to update status: %v", err)
			}
		}
	}
}
