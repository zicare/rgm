package mysql

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

//Update exported
func (Table) Update(qo *ds.QueryOptions) (int64, error) {

	t, ok := qo.DataSource.(ITable)
	if !ok {
		return 0, new(NotITableError)
	}

	tx, _ := Db().Begin()

	if err := t.BeforeUpdate(qo, tx); err != nil {
		return 0, err
	}

	b := sqlbuilder.MySQL.NewUpdateBuilder()
	b.Update(t.Name())

	assignments := []string{}
	fv := ds.Values(qo.DataSource)
	for _, f := range qo.WritableFields {
		assignments = append(assignments, b.Assign(f, fv[f]))
	}
	b.Set(assignments...)

	// set where Equal for Primary param
	for k, v := range qo.Equal[ds.Primary] {
		b.Where(b.Equal(k, v))
	}

	// set where Equal for Url param
	for k, v := range qo.Equal[ds.Url] {
		b.Where(b.Equal(k, v))
	}

	// set where EqualQPar
	for k, v := range qo.Equal[ds.Qry] {
		b.Where(b.Equal(k, v))
	}

	// set where IsNull
	for _, j := range qo.IsNull {
		b.Where(b.IsNull(j))
	}

	// set where IsNotNull
	for _, j := range qo.IsNotNull {
		b.Where(b.IsNotNull(j))
	}

	// set where In
	for k, v := range qo.In {
		b.Where(b.In(k, v...))
	}

	// set where NotIn
	for k, v := range qo.NotIn {
		b.Where(b.NotIn(k, v...))
	}

	// set where NotEqual
	for k, v := range qo.NotEqual {
		b.Where(b.NotEqual(k, v))
	}

	// set where GreaterThan
	for k, v := range qo.GreaterThan {
		b.Where(b.GreaterThan(k, v))
	}

	// set where GreaterEqualThan
	for k, v := range qo.GreaterEqualThan {
		b.Where(b.GreaterEqualThan(k, v))
	}

	// set where LessThan
	for k, v := range qo.LessThan {
		b.Where(b.LessThan(k, v))
	}

	// set where LessEqualThan
	for k, v := range qo.LessEqualThan {
		b.Where(b.LessEqualThan(k, v))
	}

	// set order by
	b.OrderBy(qo.Order...)

	// set limit
	if qo.Limit != nil {
		b.Limit(*qo.Limit)
	}

	q, args := b.Build()

	//fmt.Println(q, args)
	//return 0, new(ds.NotAllowedError)

	if res, err := tx.Exec(q, args...); err != nil {
		tx.Rollback()
		if me, ok := err.(*mysql.MySQLError); ok {
			switch me.Number {
			case 1048:
				// Column 'x' cannot be null
				s := strings.Split(me.Message, "'")
				e := ds.UpdateError{Message: msg.Get("24").SetField(s[1]).SetArgs("null", "required", "")}
				return 0, &e
			case 1451:
				// 1451 - Cannot delete or update a parent row, can happen if trying to update value of the primary key column
				return 0, new(ds.ForeignKeyConstraint)
			case 1452:
				// 1452 - Cannot add or update a child row
				return 0, new(ds.ForeignKeyConstraint)
			}
		}
		return 0, err
	} else if err := t.AfterUpdate(qo, tx); err != nil {
		tx.Rollback()
		return 0, err
	} else if rows, err := res.RowsAffected(); err != nil {
		tx.Commit()
		return 0, err
	} else {
		tx.Commit()
		return rows, nil
	}

}
