package examples_test

import (
	"context"
	"github.com/erichnascimento/golib/database/transaction"
	"github.com/erichnascimento/golib/database/transaction/sql"
)
type AddressStorage struct {
	connection sql.Connection
}

func (s *AddressStorage) Add(ctx context.Context, addr string) error {
	//if transaction.InTransaction(ctx) {
	//	transaction.Exec(ctx, "INSERT INTO addresses (street) VALUES ($1)", addr)
	//}
	//return s.stmtAdd.Exec(ctx, "INSERT INTO addresses (street) VALUES ($1)", addr)
	_, err := s.connection.ExecContext(ctx, "INSERT INTO addresses (address) VALUES ($1)", addr)
	return err
}

func (s *AddressStorage) Get(ctx context.Context, id interface{}) (interface{}, error) {
	// try execute statement
	res, err := transaction.ExecStmt(ctx, s.addStmt, addr)
	if err == "in transaction" {
		// try raw mode
		err, res = transaction.Exec(ctx, "INSERT INTO addresses (street) VALUES ($1)", addr)
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	return "av mal rondon", nil
}

type AddressService struct {
	storage *AddressStorage
}

func (s *AddressService) Add(ctx context.Context, addr interface{}) error {
	ctx = transaction.Begin(ctx) // when reading, do not do that in order to use the main transaction
	err := s.storage.Add(ctx, addr)
	if err != nil {
		transaction.Abort(ctx)
		return err
	}
	transaction.Commit(ctx)
	return nil
}

func (s *AddressService) Get(ctx context.Context, id interface{}) (interface{}, error) {
	addr, err := s.storage.Get(ctx, id)
	if err != nil {
		transaction.Abort(ctx)
		return nil, err
	}

	return addr, nil
}

type CustomerStorage struct {
	connection sql.Connection
}

func (s *CustomerStorage) Add(ctx context.Context, customer interface{}) error {
	s.sqlExecutor.Exec(ctx, "INSERT INTO customer (name) VALUES ($1)", customer)
	return nil

	//if transaction.InTransaction(ctx) {
	//	transaction.Exec(ctx, "INSERT INTO customer (name) VALUES ($1)", customer)
	//}
	//return s.stmtAdd.Exec(ctx, "INSERT INTO addresses (street) VALUES ($1)", customer)
}

type CustomerService struct {
	storage *CustomerStorage
	addressService *AddressService
}

func (s *CustomerService) Add(ctx context.Context, customer, addr string) error {
	ctx = transaction.Begin(ctx)

	err := s.storage.Add(ctx, customer)
	if err != nil {
		transaction.Abort(ctx)
		return err
	}

	err = s.addressService.Add(ctx, addr)
	if err != nil {
		transaction.Abort(ctx)
		return err
	}

	transaction.Commit(ctx)
	return nil
}