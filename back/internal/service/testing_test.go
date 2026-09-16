package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/service/mocks"
	"capuchin/internal/sqlite"
	sqlitemocks "capuchin/internal/sqlite/mocks"
)

// newTxManager возвращает мок TransactionManager.
func newTxManager(t *testing.T, ctx context.Context, err error) *mocks.TransactionManagerMock {
	t.Helper()

	txManager := mocks.NewTransactionManagerMock(t)
	txManager.
		On("ExecuteInTransaction", ctx,
			mock.AnythingOfType("func(context.Context, sqlite.SQLExecutor) error"),
			mock.AnythingOfType("func(context.Context, sqlite.SQLExecutor) error"),
		).
		Run(func(args mock.Arguments) {
			tx := sqlitemocks.NewSQLExecutorMock(t)

			fn1, ok := args.Get(1).(func(context.Context, sqlite.SQLExecutor) error)
			assert.True(t, ok)

			fn2, ok := args.Get(2).(func(context.Context, sqlite.SQLExecutor) error)
			assert.True(t, ok)

			_ = fn1(ctx, tx)
			_ = fn2(ctx, tx)
		}).
		Return(err)

	return txManager
}
