package mysql

import (
	"database/sql"
	"errors"

	"github.com/EricChiou/gooq"
)

// Gooq instance
type Gooq struct {
	SQL  gooq.SQL
	Args []interface{}
}

// AddValues add args
func (g *Gooq) AddValues(values ...interface{}) {
	g.Args = append(g.Args, values...)
}

// Exec execute sql
func (g *Gooq) Exec(sql string, args ...interface{}) (sql.Result, error) {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, errors.New("db not initialized")
	}

	return tx.Exec(sql, args...)
}

// Query query rows
func (g *Gooq) Query(sql string, rowMapper func(*sql.Rows) error, args ...interface{}) error {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return errors.New("db not initialized")
	}

	rows, err := tx.Query(sql, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	return rowMapper(rows)
}

// QueryRow query one row
func (g *Gooq) QueryRow(sql string, rowMapper func(row *sql.Row) error, args ...interface{}) error {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return errors.New("db not initialized")
	}

	row := tx.QueryRow(sql, args...)

	return rowMapper(row)
}
