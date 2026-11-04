package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/cursor"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/extractor"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/loader"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
	"github.com/milana-dev/postgres-clickhouse-etl/internal/transform"
)

func newPool(ctx context.Context, envFile string) (*pgxpool.Pool, error) {
	envMap, err := godotenv.Read(envFile)
	if err != nil {
		return nil, err
	}
	value := envMap["DATABASE_URL"]
	if value == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty in %q", envFile)
	}
	pool, err := pgxpool.New(ctx, value)
	if err != nil {
		return nil, err
	}
	return pool, err
}

func connect(ctx context.Context, envFile string) (clickhouse.Conn, error) {
	envMap, err := godotenv.Read(envFile)
	if err != nil {
		return nil, err
	}
	addr := envMap["CLICKHOUSE_ADDR"]
	if addr == "" {
		return nil, fmt.Errorf("CLICKHOUSE_ADDR is empty in %q", envFile)
	}
	pass := envMap["CLICKHOUSE_PASSWORD"]

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: pass,
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func main() {

	ctx, stopSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignal()

	poolA, err := newPool(ctx, ".env.bank_a")
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer poolA.Close()

	poolB, err := newPool(ctx, ".env.bank_b")
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer poolB.Close()

	poolC, err := newPool(ctx, ".env.bank_c")
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer poolC.Close()

	conn, err := connect(ctx, ".env.clickhouse")
	if err != nil {
		log.Fatalf("unable to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	cursorA, cursorUUIDA, err := cursor.Last(ctx, conn, "bank_a")
	if err != nil {
		log.Fatalf("unable to get cursor for bank A: %v", err)
	}

	chA := transform.NormalizeBankA(ctx, extractor.ExtractorBankA(ctx, poolA, cursorA, cursorUUIDA, 500, 1*time.Second))

	cursorB, cursorUUIDB, err := cursor.Last(ctx, conn, "bank_b")
	if err != nil {
		log.Fatalf("unable to get cursor for bank B: %v", err)
	}

	chB := transform.NormalizeBankB(ctx, extractor.ExtractorBankB(ctx, poolB, cursorB, cursorUUIDB, 500, 1*time.Second))

	cursorC, cursorUUIDC, err := cursor.Last(ctx, conn, "bank_c")
	if err != nil {
		log.Fatalf("unable to get cursor for bank c: %v", err)
	}

	var chC <-chan models.Transaction

	if cursorUUIDC == uuid.Nil {
		chC = transform.NormalizeBankC(ctx, extractor.ExtractorBankC(ctx, poolC, 0, cursorUUIDC, 500, 1*time.Second))
	} else {
		chC = transform.NormalizeBankC(ctx, extractor.ExtractorBankC(ctx, poolC, cursorC.Unix(), cursorUUIDC, 500, 1*time.Second))
	}

	res := transform.Merge(ctx, chA, chB, chC)
	batch := loader.Batch(ctx, res, 10, 5*time.Second)
	err = loader.Load(ctx, conn, batch, 3, time.Second)
	if err != nil {
		log.Printf("loader failed: %v", err)
	}
}
