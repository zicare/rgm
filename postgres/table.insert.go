package postgres

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	"github.com/zicare/rgm/ds"
)

//Insert exported
func (Table) Insert(qo *ds.QueryOptions) error {

	t, ok := qo.DataSource.(ITable)
	if !ok {
		return new(NotITableError)
	}

	tx, _ := Db().Begin()

	if err := t.BeforeInsert(qo, tx); err != nil {
		return err
	}

	b := sqlbuilder.PostgreSQL.NewInsertBuilder()
	b.InsertInto(t.Name())

	b.Cols(qo.WritableFields...)
	fv := ds.Values(qo.DataSource)
	var wvals []interface{}
	for _, f := range qo.WritableFields {
		wvals = append(wvals, fv[f])
	}
	b.Values(wvals...)

	q, args := b.BuildWithFlavor(sqlbuilder.PostgreSQL)

	//fmt.Println(q, args)
	//return new(ds.NotAllowedError)

	s := sqlbuilder.NewStruct(t).For(sqlbuilder.PostgreSQL)
	if err := tx.QueryRow(q+" RETURNING *", args...).Scan(s.AddrWithCols(qo.Fields, &t)...); err != nil {
		tx.Rollback()
		if me, ok := err.(*pq.Error); ok && me.Code == "23505" { //unique_violation
			// Duplicated entry
			return new(ds.DuplicatedEntry)
		} else if ok && me.Code == "23503" { //foreign_key_violation
			// Cannot add or update a child row
			return new(ds.ForeignKeyConstraint)
		} else if ok && me.Code == "23000" { //integrity_constraint_violation
			// Validation error
			return new(ds.ValidationError)
		}
		return err
	} else if err := t.AfterInsert(qo, tx); err != nil {
		tx.Rollback()
		return err
	}

	//return find(c, m, PID(m, fields.Primary), false)

	tx.Commit()
	return nil

}
