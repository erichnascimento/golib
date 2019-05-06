package storage

import (
	"database/sql"
	"golang.org/x/xerrors"
)

type TxReader interface {
	Tx
	Reader
}

type TxWriter interface {
	Tx
	Writer
}

type TxRW struct {
	Tx
	ReaderWriter
}

type Tx interface {
	Commit() error
	Rollback() error
	CommitOrRollback(*error) error
	CommitOrRollbackWithErrorHandler(*error, func(error))
}

type tx struct {
	*sql.Tx
}

func (tx *tx) CommitOrRollback(err *error) (resultErr error) {
	tx.CommitOrRollbackWithErrorHandler(err, func(err1 error) {
		resultErr = err1
	})
	return
}

func (tx *tx) CommitOrRollbackWithErrorHandler(err *error, errHandler func(error)) {
	if errHandler == nil {
		panic("errHandler can not be nil")
	}

	if *err != nil {
		// Log received error
		errHandler(xerrors.Errorf("rolling back by error: %w", *err))

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			errHandler(xerrors.Errorf("unable to rollback transaction: %w", rollbackErr))
		}
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		errHandler(xerrors.Errorf("unable to commit transaction: %w", commitErr))
	}
	return
}


