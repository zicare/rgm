package mysql

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

//Insert exported
func (Table) Insert(qo *ds.QueryOptions) error {

	t, ok := qo.DataSource.(ITable)
	if !ok {
		return new(NotITableError)
	}

	if err := t.BeforeInsert(qo); err != nil {
		return err
	}

	b := sqlbuilder.MySQL.NewInsertBuilder()
	b.InsertInto(t.Name())

	b.Cols(qo.WritableFields...)
	b.Values(qo.WritableValues...)

	q, args := b.Build()

	s := sqlbuilder.NewStruct(t).For(sqlbuilder.MySQL)
	if err := Db().QueryRow(q+" RETURNING *", args...).
		Scan(s.AddrWithCols(qo.Fields, &t)...); err != nil {
		return err
	}

	//return find(c, m, PID(m, fields.Primary), false)

	return nil

}
