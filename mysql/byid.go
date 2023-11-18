package mysql

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

func ByID(tx *sql.Tx, t ITable, v ...interface{}) error {

	k, _, _, err := ds.Meta(t)
	if err != nil {
		return err
	}

	s := sqlbuilder.NewStruct(t)
	b := s.SelectFrom(t.Name())
	for inx, key := range k {
		b.Where(b.Equal(key, v[inx]))
	}

	q, args := b.Build()
	if err := tx.QueryRow(q, args...).Scan(s.Addr(&t)...); err == sql.ErrNoRows {
		return new(ds.NotFoundError)
	}
	return nil
}
