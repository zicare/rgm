package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

// Queryable lets this work with *sql.DB and *sql.Tx.
type Queryable interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// ByID is a convenience wrapper around ByIDX using context.Background() and mysql.Db().
// Usage: _ = ByID(myModel, 123)
func ByID(m ds.IDataSource, pkVals ...any) error {
	return ByIDX(context.Background(), Db(), m, pkVals...)
}

// ByIDX can run with either a *sql.DB or *sql.Tx.
// Usage with tx: _ = ByIDX(ctx, tx, myModel, 123)
// Usage without tx: _ = ByIDX(ctx, Db(), myModel, 123)
func ByIDX(ctx context.Context, q Queryable, m ds.IDataSource, pkVals ...any) error {
	keys, _, _, _, err := ds.Meta(m)
	if err != nil {
		return err
	}
	if len(pkVals) != len(keys) {
		return errors.New("ByID[X]: pk arity mismatch")
	}

	s := sqlbuilder.NewStruct(m)
	b := s.SelectFrom(m.Name())
	b.Select(s.Columns()...) // ensure stable column order for Scan
	for i, col := range keys {
		b.Where(b.Equal(col, pkVals[i]))
	}
	sqlStr, args := b.Build()

	if err := q.QueryRowContext(ctx, sqlStr, args...).Scan(s.Addr(&m)...); err != nil {
		if err == sql.ErrNoRows {
			return new(ds.NotFoundError)
		}
		return err
	}
	return nil
}
