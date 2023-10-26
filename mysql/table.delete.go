package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

// Delete supports single and multiple records removal.
// It first checks with Table's BeforeDelete method for extra constraints.
// BeforeDelete can also return a *ds.NotAllowedError to abort Delete.
// Beware that qo.DataSource must implement ds.IDataSource.
func (Table) Delete(qo *ds.QueryOptions) (int64, error) {

	t, ok := qo.DataSource.(ds.IDataSource)
	if !ok {
		return 0, new(ds.NotIDataSourceError)
	}

	b := sqlbuilder.DeleteFrom(t.Name())

	// BeforeDelete check
	if where, err := t.BeforeDelete(qo); err != nil {
		return 0, err
	} else {
		// set where scope
		for k, v := range where {
			b.Where(b.Equal(k, v))
		}
	}

	// set where EqualPk
	for k, v := range qo.Equal[ds.Primary] {
		b.Where(b.Equal(k, v))
	}

	// set where EqualUPar
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

	// build the sql
	q, args := b.Build()

	//fmt.Println(q, args)
	//return 0, nil

	// Execute delete
	if res, err := Db().Exec(q, args...); err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1451 {
			// Cannot delete or update a parent row
			return 0, new(ds.ForeignKeyConstraint)
		}
		return 0, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return 0, err
	} else {
		return rows, nil
	}
}
