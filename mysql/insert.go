package mysql

import (
	"github.com/go-sql-driver/mysql"
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
	fv := ds.Values(qo.DataSource)
	var wvals []interface{}
	for _, f := range qo.WritableFields {
		wvals = append(wvals, fv[f])
	}
	b.Values(wvals...)

	q, args := b.Build()

	//fmt.Println(q, args)
	//return new(ds.NotAllowedError)

	s := sqlbuilder.NewStruct(t).For(sqlbuilder.MySQL)
	if err := Db().QueryRow(q+" RETURNING *", args...).Scan(s.AddrWithCols(qo.Fields, &t)...); err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1062 {
			// Duplicated entry
			return new(ds.DuplicatedEntry)
		} else if ok && me.Number == 1452 {
			// Cannot add or update a child row
			return new(ds.ForeignKeyConstraint)
		}

		return err
	}

	//return find(c, m, PID(m, fields.Primary), false)

	return nil

}
