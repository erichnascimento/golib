package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/erichnascimento/golib/database/transaction"
)

var dbs = map[string]*sql.DB{}
var connectors = map[string]Connection{}
var registeredTransactions = map[interface{}]Transaction{}

func RegisterConnector(name string, db *sql.DB) {
	dbs[name] = db
}

func Connect(name string) Connection {
	for cname, c :=  range connectors {
		if cname == name {
			return c
		}
	}

	db := dbs[name]
	c := &connection{
		db:db,
	}
	connectors[name] = c
	return c
}

type Transaction interface {
	QueryRowContext(ctx context.Context, query string, args... interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args... interface{}) (sql.Result, error)
	Rollback() error
	Commit() error
}

type sqlTransaction struct {
	*sql.Tx
}

type sqlAutoCommittedTransaction struct {
	*sql.DB
}

func (tx *sqlAutoCommittedTransaction) Rollback() error {
	return errors.New("rollback not supported for auto committed transaction")
}

func (tx *sqlAutoCommittedTransaction) Commit() error {
	return errors.New("rollback not supported for auto committed transaction")
}

type Connection interface {
	// Put here sql.DB methods useful for storages
	ExecContext(ctx context.Context, query string, args... interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error)
}

type connection struct {
	db *sql.DB
}

func (c *connection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, err := c.resolveTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx.ExecContext(ctx, query, args...)
}

func (c *connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	tx, err := c.resolveTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx.QueryRowContext(ctx, query, args...), nil
}

func (c *connection) resolveTx(ctx context.Context) (Transaction, error) {
	id := transaction.Attach(ctx, c)
	if tx, ok := registeredTransactions[id]; ok {
		return tx, nil
	}

	var tx Transaction
	if transaction.IsMainTransaction(id) {
		tx = &sqlAutoCommittedTransaction{
			DB: c.db,
		}
	} else {
		dbTx, err := c.db.Begin()
		if err != nil {
			return nil, err
		}
		tx = &sqlTransaction{
			Tx: dbTx,
		}
	}
	registeredTransactions[id] = tx

	return tx, nil
}

func (c *connection) Commit(id interface{}) error {
	return registeredTransactions[id].Commit()
}

func (c *connection) Abort(id interface{}) error {
	return registeredTransactions[id].Rollback()
}
