package qgenda

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Staffs []Staff

func (s *Staffs) Get(ctx context.Context, c *Client, rc *RequestConfig) error {
	return nil
}

func (s *Staffs) Process() error {
	return nil
}

func (s *Staffs) LoadFile(filename string) error {
	return nil

}

func (s *Staffs) PGInsertRows(ctx context.Context, tx *sqlx.Tx, schema, tablename, id string) (sql.Result, error) {
	return nil, nil
}

func (s *Staffs) EPL(ctx context.Context, c *Client, rc *RequestConfig,
	db *sqlx.DB, schema, table string, newRowsOnly bool) (sql.Result, error) {
	return nil, nil
}
