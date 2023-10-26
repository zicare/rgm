package mysql

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

// Count returns the number of qo.DataSource records that match qo settings.
// Beware that qo.DataSource must implement ds.IDataSource.
func (Table) Count(qo *ds.QueryOptions) (count int64, err error) {

	s := sqlbuilder.NewStruct(qo.DataSource)
	b := s.SelectFrom(qo.DataSource.Name())

	// set where Equal for Url param
	for k, v := range qo.Equal[ds.Url] {
		b.Where(b.Equal(k, v))
	}

	// set where Equal for Query param
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

	// get total count
	b.Select(b.As("COUNT(*)", "t"))

	q, args := b.Build()

	if err := Db().QueryRow(q, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil

}
