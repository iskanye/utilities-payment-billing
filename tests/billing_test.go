package tests

import (
	"io"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/utilities-payment-billing/tests/suite"
	"github.com/iskanye/utilities-payment-proto/billing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	deltaDay = 60 * 60 * 24
)

func amount() int32 {
	return int32(gofakeit.Number(0, 1000000))
}

func userID() int64 {
	return int64(gofakeit.Number(1, 100000))
}

func CheckDueDate(t *testing.T, s *suite.Suite, dueDate string) {
	term, err := time.Parse(time.RFC3339, dueDate)
	require.NoError(t, err)

	assert.InDelta(t, time.Now().AddDate(0, s.Cfg.Term, 0).Unix(), term.Unix(), deltaDay)
}

func TestCreateBill_Success(t *testing.T) {
	ctx, s := suite.New(t)

	address := gofakeit.Address().Address
	amount := amount()
	userID := userID()

	respBill, err := s.BillingClient.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  amount,
		UserId:  userID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respBill)

	respBills, err := s.BillingClient.GetBills(ctx, &billing.BillsRequest{
		UserId: userID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respBill)

	for {
		bill, err := respBills.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.NotEmpty(t, bill)

		assert.Equal(t, respBill.GetBillId(), bill.GetBillId())
		assert.Equal(t, address, bill.GetAddress())
		assert.Equal(t, amount, bill.GetAmount())
		assert.Equal(t, userID, bill.GetUserId())

		CheckDueDate(t, s, bill.GetDueDate())
	}
}

func TestCreateBill_Dublicates(t *testing.T) {
	ctx, s := suite.New(t)

	address := gofakeit.Address().Address
	amount := amount()
	userID := userID()
	var ids []int64

	respBill, err := s.BillingClient.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  amount,
		UserId:  userID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respBill)
	ids = append(ids, respBill.GetBillId())

	respBill, err = s.BillingClient.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  amount,
		UserId:  userID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respBill)
	ids = append(ids, respBill.GetBillId())

	respBills, err := s.BillingClient.GetBills(ctx, &billing.BillsRequest{
		UserId: userID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respBill)

	var bills []*billing.Bill

	for {
		bill, err := respBills.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.NotEmpty(t, bill)

		bills = append(bills, bill)
	}

	assert.Equal(t, 2, len(bills))

	for i, bill := range bills {
		assert.Equal(t, ids[i], bill.GetBillId())
		assert.Equal(t, address, bill.GetAddress())
		assert.Equal(t, amount, bill.GetAmount())
		assert.Equal(t, userID, bill.GetUserId())

		CheckDueDate(t, s, bill.GetDueDate())
	}
}

func TestCreateBill_FailCases(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		address     string
		amount      int32
		userID      int64
		expectedErr string
	}{
		{
			name:        "CreateBill with no address",
			address:     "",
			amount:      amount(),
			userID:      userID(),
			expectedErr: "address is required",
		},
		{
			name:        "CreateBill with zero amount",
			address:     gofakeit.Address().Address,
			amount:      0,
			userID:      userID(),
			expectedErr: "amount must be positive",
		},
		{
			name:        "CreateBill with negative amount",
			address:     gofakeit.Address().Address,
			amount:      -1000,
			userID:      userID(),
			expectedErr: "amount must be positive",
		},
		{
			name:        "CreateBill with no userID",
			address:     gofakeit.Address().Address,
			amount:      amount(),
			userID:      0,
			expectedErr: "user_id is required",
		},
		{
			name:        "CreateBill with no data",
			address:     "",
			amount:      0,
			userID:      0,
			expectedErr: "address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.BillingClient.AddBill(ctx, &billing.Bill{
				Address: tt.address,
				Amount:  tt.amount,
				UserId:  tt.userID,
			})

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
