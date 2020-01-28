package transaction

import (
	"context"
	"errors"
	"github.com/erichnascimento/golib/uuid"
)

const txIDKeyName = "golib.tx.id"
const txLvlKeyName = "golib.tx.lvl"

type level int

func (l level) String() string {
	switch l {
	case levelRoot:
		return "root"
	case levelChild:
		return "child"
	default:
		return ""
	}
}

const (
	levelRoot level = iota
	levelChild

	mainTransactionID = "main"
)

type Transaction interface {
	Commit(id interface{}) error
	Abort(id interface{}) error
}

var activeTransactions = map[string][]Transaction{}

func Begin(ctx context.Context) context.Context {
	if ctx == nil {
		panic("ctx can not be nil")
	}

	lvl := levelChild
	txID, ok := ctx.Value(txIDKeyName).(string)
	if !ok {
		txID = uuid.NewUUID().String()
		lvl = levelRoot
		activeTransactions[txID] = []Transaction{}
	}
	ctx = context.WithValue(ctx, txIDKeyName, txID)
	return context.WithValue(ctx, txLvlKeyName, lvl)
}

func IsMainTransaction(id string) bool {
	return id == mainTransactionID
}

func Attach(ctx context.Context, tx Transaction)string {
	id, ok := ctx.Value(txIDKeyName).(string)
	if !ok {
		id = mainTransactionID
	}
	if !transactionExists(id, tx) {
		if transactions, ok := activeTransactions[id]; !ok {
			activeTransactions[id] = []Transaction{tx}
		} else {
			activeTransactions[id] = append(transactions, tx)
		}
	}

	return id
}

func transactionExists(id string, tx Transaction) bool {
	transactions, ok := activeTransactions[id]
	if !ok {
		return false
	}

	for _, transaction := range transactions {
		if transaction == tx {
			return true
		}
	}

	return false
}

func Commit(ctx context.Context) (Result, error) {
	r := newResult(ctx)
	transactions, err := forEachTransaction(r.id)
	if err != nil {
		return nil, errors.New("no transaction for given context")
	}

	for tx := range transactions {
		r.n++
		if err := tx.Commit(r.id); err != nil {
			r.addErr(err)
		}
	}

	return r, nil
}

func forEachTransaction(id string) (<-chan Transaction, error) {
	transactions, ok := activeTransactions[id]
	if !ok {
		return nil, errors.New("no transaction for given context")
	}

	transactionsCh := make(chan Transaction)
	go func() {
		defer close(transactionsCh)
		for _, tx := range transactions {
			transactionsCh <- tx
		}
	}()

	return transactionsCh, nil
}

type Result interface {
	ID() interface{}
	Errors() []error
	N() int
}

func newResult(ctx context.Context) *result {
	id, ok :=  ctx.Value(txIDKeyName).(string)
	if !ok {
		id = mainTransactionID
	}
	return &result{id: id}
}

type result struct {
	id string
	errors []error
	n int
}

func (r *result) ID() interface{} {
	return r.id
}

func (r *result) Errors() []error {
	return r.errors
}

func (r *result) addErr(err error) {
	if r.errors == nil {
		r.errors = make([]error, 0)
	}
	r.errors = append(r.errors, err)
}

func (r *result) N() int {
	return r.n
}