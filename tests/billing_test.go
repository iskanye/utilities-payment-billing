package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/utilities-payment-billing/tests/suite"
	"github.com/iskanye/utilities-payment-proto/billing"
	"github.com/stretchr/testify/assert"
)

func TestCreateBill(t *testing.T) {
	ctx, s := suite.New(t)

	address := gofakeit.Address().Address
	amount := gofakeit.Number(0, 1000000)

	respBill, err := s.BillingClient.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  int32(amount),
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, respBill)
}
