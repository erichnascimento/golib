package storage

import (
	"database/sql"
	"golang.org/x/xerrors"
)

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

type TransactionManager struct {
	db *sql.DB
}

func (tm *TransactionManager) Begin() (TxRW, error)  {
	dbTx, err := tm.db.Begin()
	if err != nil {
		return nil, xerrors.Errorf("begin transaction error: %w", err)
	}

	tx := &tx{dbTx}

	return &txRW{
		TxRW: tx,
	}, nil
}
