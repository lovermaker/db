package db

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Interface client interface
type Interface[T any] interface {
	StartTransaction(opts ...*sql.TxOptions) TransactionInterface[T]
	DAOProvider[T]
}

// DAOProvider provide db connection
type DAOProvider[T any] interface {
	DAO() T
}

// Client DAO manager, provider sets of db operations
type Client[T any] struct {
	StartTransactionFunc func(opts ...*sql.TxOptions) TransactionInterface[T]
	DAOFunc              func() T
}

func (c *Client[T]) DAO() T {
	return c.DAOFunc()
}

func NewClient[T any](creator func(db *gorm.DB) T, config *Config) (Interface[T], error) {
	return injectClient(&Client[T]{}, creator, config)
}

func injectClient[T any](client *Client[T], creator func(db *gorm.DB) T, config *Config) (*Client[T], error) {
	impl, err := newClient(creator, config)
	if err != nil {
		return nil, err
	}
	client.StartTransactionFunc = impl.StartTransaction
	return client, nil
}

func NewMockClient[T any](config *Config) *Client[T] {
	return &Client[T]{}
}

func (c *Client[T]) StartTransaction(opts ...*sql.TxOptions) TransactionInterface[T] {
	return c.StartTransactionFunc(opts...)
}

type client[T any] struct {
	db      *gorm.DB
	dao     T
	creator func(db *gorm.DB) T
}

func newClient[T any](creator func(db *gorm.DB) T, config *Config) (*client[T], error) {
	if config.Charset == "" {
		config.Charset = "utf8"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.UserName,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("init database error: %s", dsn)
	}
	if config.Debug {
		db.Debug()
	}
	cli := &client[T]{
		db:      db,
		creator: creator,
		dao:     creator(db),
	}
	return cli, nil
}

// StartTransaction start a transaction
func (c *client[T]) StartTransaction(opts ...*sql.TxOptions) TransactionInterface[T] {
	tx := c.db.Begin(opts...)
	return newTransactionClient[T](c.creator, tx)
}

// DAO returns dao object
func (c *client[T]) DAO() T {
	return c.dao
}
