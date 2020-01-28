package sql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erichnascimento/golib/database/transaction"
	"testing"
)

const connectionName = "db-for-name-and-age"

func Test(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	RegisterConnector(connectionName, db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO names").WithArgs("edu").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT name FROM names WHERE name").WithArgs("edu").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("edu"))
	mock.ExpectCommit()

	nameStorage := NewNameStorage()
	ctx := context.Background()
	ctx = transaction.Begin(ctx)
	ctx = transaction.Begin(ctx)
	expectedName := "edu"
	err = nameStorage.Add(ctx, expectedName)
	if err != nil {
		t.Fatal(err)
	}

	givenName, err := nameStorage.Get(ctx, "edu")
	if err != nil {
		t.Fatal(err)
	}
	if givenName != "edu" {
		t.Fatalf("expected=%s, given=%s", expectedName, givenName)
	}
	res, err := transaction.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if res.N() != 1 {
		t.Fatalf("expected=%d, given=%d", 1, res.N())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

type NameStorage struct {
	connection Connection
}

func NewNameStorage() *NameStorage {
	return &NameStorage{
		connection:Connect(connectionName),
	}
}

func (s *NameStorage) Add(ctx context.Context, name string) error {
	_, err := s.connection.ExecContext(ctx, "INSERT INTO names (name) VALUES ($1)", name)
	if err != nil {
		return err
	}
	return nil
}

func (s *NameStorage) Get(ctx context.Context, id string) (name string, err error) {
	r, err := s.connection.QueryRowContext(ctx, "SELECT name FROM names WHERE name = $1)", id)
	if err != nil {
		return
	}
	err = r.Scan(&name)
	return
}
