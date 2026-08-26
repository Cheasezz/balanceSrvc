package tests

import (
	"context"
	"testing"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	testsuite "github.com/Cheasezz/balanceSrvc/tests/suite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcBalance_SystemTransactionTo(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)

	tests := []struct {
		name    string
		req     *blnc.SystemTrxToRequest
		wantErr error
	}{
		{
			name: "happy path",
			req: &blnc.SystemTrxToRequest{
				IdempotencyKey: uuid.NewString(),
				UserId:         uuid.NewString(),
				SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:         10000,
			},
			wantErr: nil,
		},
		{
			name: "error bad userId",
			req: &blnc.SystemTrxToRequest{
				IdempotencyKey: uuid.NewString(),
				UserId:         "bad uuid",
				SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:         10000,
			},
			wantErr: status.Error(codes.InvalidArgument, core.ErrInvalidUuid.Error()),
		},
		{
			name: "error zero amount",
			req: &blnc.SystemTrxToRequest{
				IdempotencyKey: uuid.NewString(),
				UserId:         uuid.NewString(),
				SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:         0,
			},
			wantErr: status.Error(codes.InvalidArgument, core.ErrInvalidAmount.Error()),
		},
		{
			name: "error invalid transaction type",
			req: &blnc.SystemTrxToRequest{
				IdempotencyKey: uuid.NewString(),
				UserId:         uuid.NewString(),
				SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
				Amount:         10000,
			},
			wantErr: status.Error(codes.InvalidArgument, core.ErrUnknownTrxType.Error()),
		},
		{
			name: "error bad idempotencyKey",
			req: &blnc.SystemTrxToRequest{
				IdempotencyKey: "bad idempotency_key",
				UserId:         uuid.NewString(),
				SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:         10000,
			},
			wantErr: status.Error(codes.InvalidArgument, core.ErrInvalidIdempotencyKey.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
			defer cancelCtx()

			_, err := suit.BalanceClient.SystemTransactionTo(ctx, tt.req)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGrpcBalance_ConcurrentSystemToUserDifferentIdempKeys(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	wantBalance := 500
	balanceIncPerTrx := 10
	user := uuid.NewString()

	errors := runConcurrent(tests, func() error {
		req := &blnc.SystemTrxToRequest{
			IdempotencyKey: uuid.NewString(),
			UserId:         user,
			SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			Amount:         uint64(balanceIncPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err := suit.BalanceClient.SystemTransactionTo(ctx, req)

		return err
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user})

	require.Nil(t, err)
	require.Exactly(t, int(currentBalance.GetBalance()), wantBalance)
	require.Len(t, errors, 0)
}

func TestGrpcBalance_ConcurrentSystemToUserSameIdempKeys(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	wantBalance := 10
	balanceIncPerTrx := 10
	user := uuid.NewString()
	iKey := uuid.NewString()

	errors := runConcurrent(tests, func() error {
		req := &blnc.SystemTrxToRequest{
			IdempotencyKey: iKey,
			UserId:         user,
			SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			Amount:         uint64(balanceIncPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err := suit.BalanceClient.SystemTransactionTo(ctx, req)

		return err
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user})

	require.Nil(t, err)
	require.Exactly(t, int(currentBalance.GetBalance()), wantBalance)
	require.Len(t, errors, tests-1)
}

func TestGrpcBalance_ConcurrentSystemToUserSameIdempKeyPerUser(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	wantBalance := 10
	balanceIncPerTrx := 10
	user1 := uuid.NewString()
	iKey1 := uuid.NewString()
	user2 := uuid.NewString()
	iKey2 := uuid.NewString()

	errors := runConcurrent(tests, func() error {
		req1 := &blnc.SystemTrxToRequest{
			IdempotencyKey: iKey1,
			UserId:         user1,
			SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			Amount:         uint64(balanceIncPerTrx),
		}

		req2 := &blnc.SystemTrxToRequest{
			IdempotencyKey: iKey2,
			UserId:         user2,
			SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			Amount:         uint64(balanceIncPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err1 := suit.BalanceClient.SystemTransactionTo(ctx, req1)
		_, err2 := suit.BalanceClient.SystemTransactionTo(ctx, req2)

		if err1 != nil {
			return err1
		}
		return err2
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance1, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user1})
	require.Nil(t, err)

	currentBalance2, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user2})
	require.Nil(t, err)

	require.Exactly(t, int(currentBalance1.GetBalance()), wantBalance)
	require.Exactly(t, int(currentBalance2.GetBalance()), wantBalance)
	require.Len(t, errors, tests-1)
}

func TestGrpcBalance_ConcurrentUserToSystemDifferentIdempKeys(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	startBalance := 1000
	wantBalance := 500
	balanceDecPerTrx := 10
	user := uuid.NewString()
	depReq := &blnc.SystemTrxToRequest{
		UserId:         user,
		SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
		Amount:         uint64(startBalance),
		IdempotencyKey: uuid.NewString(),
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	_, err := suit.BalanceClient.SystemTransactionTo(ctx, depReq)

	require.NoError(t, err)

	errors := runConcurrent(tests, func() error {
		req := &blnc.SystemTrxFromRequest{
			IdempotencyKey: uuid.NewString(),
			UserId:         user,
			SystemTrxType:  blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			Amount:         uint64(balanceDecPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err := suit.BalanceClient.SystemTransactionFrom(ctx, req)

		return err
	})

	ctx, cancelCtx = context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user})

	require.Nil(t, err)
	require.Exactly(t, int(currentBalance.GetBalance()), wantBalance)
	require.Len(t, errors, 0)
}

func TestGrpcBalance_ConcurrentUserToSystemSameIdempKeys(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	startBalance := 1000
	wantBalance := 990
	balanceDecPerTrx := 10
	iKey := uuid.NewString()
	user := uuid.NewString()
	depReq := &blnc.SystemTrxToRequest{
		UserId:         user,
		SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
		Amount:         uint64(startBalance),
		IdempotencyKey: uuid.NewString(),
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	_, err := suit.BalanceClient.SystemTransactionTo(ctx, depReq)

	require.NoError(t, err)

	errors := runConcurrent(tests, func() error {
		req := &blnc.SystemTrxFromRequest{
			IdempotencyKey: iKey,
			UserId:         user,
			SystemTrxType:  blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			Amount:         uint64(balanceDecPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err := suit.BalanceClient.SystemTransactionFrom(ctx, req)

		return err
	})

	ctx, cancelCtx = context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user})

	require.Nil(t, err)
	require.Exactly(t, int(currentBalance.GetBalance()), wantBalance)
	require.Len(t, errors, tests-1)
}

func TestGrpcBalance_ConcurrentUserToSystemSameIdempKeyPerUser(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)
	tests := 50
	startBalance := 500
	wantBalance := 490
	balanceDecPerTrx := 10
	user1 := uuid.NewString()
	iKey1 := uuid.NewString()
	user2 := uuid.NewString()
	iKey2 := uuid.NewString()

	depReq1 := &blnc.SystemTrxToRequest{
		UserId:         user1,
		SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
		Amount:         uint64(startBalance),
		IdempotencyKey: uuid.NewString(),
	}

	depReq2 := &blnc.SystemTrxToRequest{
		UserId:         user2,
		SystemTrxType:  blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
		Amount:         uint64(startBalance),
		IdempotencyKey: uuid.NewString(),
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	_, err := suit.BalanceClient.SystemTransactionTo(ctx, depReq1)
	require.NoError(t, err)

	_, err = suit.BalanceClient.SystemTransactionTo(ctx, depReq2)
	require.NoError(t, err)

	errors := runConcurrent(tests, func() error {
		req1 := &blnc.SystemTrxFromRequest{
			IdempotencyKey: iKey1,
			UserId:         user1,
			SystemTrxType:  blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			Amount:         uint64(balanceDecPerTrx),
		}

		req2 := &blnc.SystemTrxFromRequest{
			IdempotencyKey: iKey2,
			UserId:         user2,
			SystemTrxType:  blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			Amount:         uint64(balanceDecPerTrx),
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
		defer cancelCtx()

		_, err1 := suit.BalanceClient.SystemTransactionFrom(ctx, req1)
		_, err2 := suit.BalanceClient.SystemTransactionFrom(ctx, req2)

		if err1 != nil {
			return err1
		}
		return err2
	})

	ctx, cancelCtx = context.WithTimeout(context.Background(), suit.CtxTimeout)
	defer cancelCtx()

	currentBalance1, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user1})
	require.Nil(t, err)

	currentBalance2, err := suit.BalanceClient.UserBalance(ctx, &blnc.BalanceRequest{UserId: user2})
	require.Nil(t, err)

	require.Exactly(t, int(currentBalance1.GetBalance()), wantBalance)
	require.Exactly(t, int(currentBalance2.GetBalance()), wantBalance)
	require.Len(t, errors, tests-1)
}
