package db

import (
	"gorm.io/gorm"
)

// TransactionInterface trx manager
type TransactionInterface[T any] interface {
	Rollback() error
	Commit() error
	DAOProvider[T]
}

// TransactionClient transaction manager
type TransactionClient[T any] struct {
	db           *gorm.DB
	CommitFunc   func() error
	RollbackFunc func() error
	DAOFunc      func() T
}

func newTransactionClient[T any](creator func(db *gorm.DB) T, tx *gorm.DB) *TransactionClient[T] {
	return injectTransactionClient[T](&TransactionClient[T]{}, creator, tx)
}

func injectTransactionClient[T any](client *TransactionClient[T], creator func(db *gorm.DB) T,
	tx *gorm.DB) *TransactionClient[T] {
	impl := &transactionClient[T]{tx: tx, dao: creator(tx)}
	client.CommitFunc = impl.Commit
	client.RollbackFunc = impl.Rollback
	client.DAOFunc = impl.DAO
	return client
}

// Rollback transaction
func (t *TransactionClient[T]) Rollback() error {
	return t.RollbackFunc()
}

// Commit transaction
func (t *TransactionClient[T]) Commit() error {
	return t.CommitFunc()
}

// DAO returns dao object
func (t *TransactionClient[T]) DAO() T {
	return t.DAOFunc()
}

type transactionClient[T any] struct {
	tx  *gorm.DB
	dao T
}

// Rollback transaction
func (t *transactionClient[T]) Rollback() error {
	return t.tx.Rollback().Error
}

// Commit transaction
func (t *transactionClient[T]) Commit() error {
	return t.tx.Commit().Error
}

// DAO return dao object
func (t *transactionClient[T]) DAO() T {
	return t.dao
}
