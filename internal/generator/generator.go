package generator

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/milana-dev/postgres-clickhouse-etl/internal/models"
	"github.com/shopspring/decimal"
)

func accountGenerator() string {

	var builder strings.Builder

	for i := 0; i < 20; i++ {
		digit := rand.IntN(10)
		digitStr := strconv.Itoa(digit)
		builder.WriteString(digitStr)

	}
	return builder.String()
}

func amountGenerator() decimal.Decimal {

	min := 1000
	max := 100000000
	randAmount := rand.IntN(max-min+1) + min
	amount := decimal.New(int64(randAmount), -2)
	return amount
}

func currencyGenerator() models.Currency {
	switch rand.IntN(3) {
	case 0:
		return models.RUB
	case 1:
		return models.EUR
	case 2:
		return models.USD
	}
	return models.RUB
}

func statusGenerator() models.Status {
	n := rand.IntN(100)
	switch {
	case n < 50:
		return models.StatusSuccess
	case n < 90:
		return models.StatusPending
	default:
		return models.StatusFailed
	}
}

func NewTransaction(bankID string) models.Transaction {

	t := models.Transaction{
		UUID:            uuid.New(),
		BankID:          bankID,
		SenderAccount:   accountGenerator(),
		ReceiverAccount: accountGenerator(),
		Amount:          amountGenerator(),
		Currency:        currencyGenerator(),
		Timestamp:       time.Now().UTC(),
		Status:          statusGenerator(),
	}
	return t
}
