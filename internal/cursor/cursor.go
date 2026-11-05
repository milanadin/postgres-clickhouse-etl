package cursor

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func Last(ctx context.Context, conn clickhouse.Conn, bankID string) (time.Time, uuid.UUID, error) {
	query := `SELECT uuid, update_at
						FROM transactions
						WHERE bank_id = ?
						ORDER BY update_at DESC, uuid DESC
						LIMIT 1`
	rows, err := conn.Query(ctx, query, bankID)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	defer rows.Close()

	var (
		cursorUUID uuid.UUID
		cursor     time.Time
	)

	if rows.Next() {
		if err := rows.Scan(&cursorUUID, &cursor); err != nil {
			return time.Time{}, uuid.Nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return time.Time{}, uuid.Nil, err
	}
	return cursor, cursorUUID, nil
}
