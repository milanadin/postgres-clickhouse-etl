package transform

import (
	"context"
	"sync"

	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
)

func Merge(ctx context.Context, inputChannels ...<-chan models.Transaction) <-chan models.Transaction {
	out := make(chan models.Transaction)
	var wg sync.WaitGroup
	for _, ch := range inputChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-ch:
					if ok {
						select {
						case out <- t:
						case <-ctx.Done():
							return

						}
					} else {
						return
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
