package database

import (
	"database/sql"
	"golang.org/x/xerrors"
	"strings"
)

type Inspector interface {
	TableExists(table string) (bool, error)
}

func NewInspector(db *sql.DB) (Inspector, error) {
	return &PostgresInspector{db}, nil
}

type PostgresInspector struct {
	db *sql.DB
}

func (pg *PostgresInspector) TableExists(table string) (bool, error) {
	const query = `
SELECT EXISTS (
   SELECT true
   FROM   information_schema.tables
   WHERE  table_schema = $1
   AND    table_name = $2
)
`
	nameParts := strings.Split(table, ".")
	r := pg.db.QueryRow(query, nameParts[0], nameParts[1])
	var exists bool
	if err := r.Scan(&exists); err != nil {
		return false, xerrors.Errorf("error querying database: %w", err)
	}
	return exists, nil
}
