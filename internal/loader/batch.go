package loader

import (
	"context"
	"time"

	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
)

func Batch(ctx context.Context, in <-chan models.Transaction, size int, timeout time.Duration) <-chan []models.Transaction {
	out := make(chan []models.Transaction)
	s := make([]models.Transaction, 0, size)
	go func() {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		defer close(out)
		for {
			select {
			case e, ok := <-in:
				if !ok {
					if len(s) > 0 {
						select {
						case out <- s:
						case <-ctx.Done():
							return
						}
					}
					return
				}
				s = append(s, e)
				if len(s) == size {
					select {
					case out <- s:
					case <-ctx.Done():
						return
					}
					s = make([]models.Transaction, 0, size)
					ticker.Reset(timeout)
				}
			case <-ticker.C:
				if len(s) > 0 {
					select {
					case out <- s:
					case <-ctx.Done():
						return
					}
					s = make([]models.Transaction, 0, size)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
