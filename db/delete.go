package db

import (
	"github.com/huandu/go-sqlbuilder"
)

// Supports single and multiple records deletion.
// It first checks with Table's BeforeDelete method
// for extra constraints that might be impossed by
// the requesting user's scope, the parent Tables, etc.
// BeforeDelete can also return a flag to abort the delete
// request, in such case Delete will return a *NotAllowedError.
func Delete(qo *QueryOptions) (int64, error) {

	dbr := sqlbuilder.DeleteFrom(qo.Table.Name())

	// BeforeDelete check
	if where, ok := qo.Table.BeforeDelete(qo.UID, qo.Parents...); !ok {
		return 0, new(NotAllowedError)
	} else {
		// set where scope
		for k, v := range where {
			dbr.Where(dbr.Equal(k, v))
		}
	}

	// set where EqualPk
	for k, v := range qo.Equal[Primary] {
		dbr.Where(dbr.Equal(k, v))
	}

	// set where EqualUPar
	for k, v := range qo.Equal[Url] {
		dbr.Where(dbr.Equal(k, v))
	}

	// set where EqualQPar
	for k, v := range qo.Equal[Query] {
		dbr.Where(dbr.Equal(k, v))
	}

	// set where IsNull
	for _, j := range qo.IsNull {
		dbr.Where(dbr.IsNull(j))
	}

	// set where IsNotNull
	for _, j := range qo.IsNotNull {
		dbr.Where(dbr.IsNotNull(j))
	}

	// set where In
	for k, v := range qo.In {
		dbr.Where(dbr.In(k, v...))
	}

	// set where NotIn
	for k, v := range qo.NotIn {
		dbr.Where(dbr.NotIn(k, v...))
	}

	// set where NotEqual
	for k, v := range qo.NotEqual {
		dbr.Where(dbr.NotEqual(k, v))
	}

	// set where GreaterThan
	for k, v := range qo.GreaterThan {
		dbr.Where(dbr.GreaterThan(k, v))
	}

	// set where GreaterEqualThan
	for k, v := range qo.GreaterEqualThan {
		dbr.Where(dbr.GreaterEqualThan(k, v))
	}

	// set where LessThan
	for k, v := range qo.LessThan {
		dbr.Where(dbr.LessThan(k, v))
	}

	// set where LessEqualThan
	for k, v := range qo.LessEqualThan {
		dbr.Where(dbr.LessEqualThan(k, v))
	}

	// set order by
	dbr.OrderBy(qo.Order...)

	// set limit
	if qo.Limit != nil {
		dbr.Limit(*qo.Limit)
	}

	// build the sql
	q, args := dbr.Build()

	//fmt.Println(q, args)
	//return 0, nil

	// Execute delete
	if res, err := Db().Exec(q, args); err != nil {
		return 0, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return 0, err
	} else {
		return rows, nil
	}
}
