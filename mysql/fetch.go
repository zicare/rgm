package mysql

import (
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
)

// Fetch returns the qo.DataStore records that match qo settings.
// Supports BeforeSelect(qo) and parent data retrieval through dig params.
// If a parent resource is not found, Fetch is aborted with a NotFoundError.
// Beware that qo.DataStore must implement ITable.
func (Table) Fetch(qo *ds.QueryOptions) (meta ds.ResultSetMeta, data []interface{}, err error) {

	t, ok := qo.DataStore.(ITable)
	if !ok {
		return meta, data, new(NotITableError)
	}

	s := sqlbuilder.NewStruct(qo.DataStore)
	b := s.SelectFrom(qo.DataStore.Name())

	// set before select constraints
	if params, err := t.BeforeSelect(qo); err != nil {
		return meta, data, new(ds.NotAllowedError)
	} else {
		for k, v := range params {
			b.Where(b.Equal(k, v))
		}
	}

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
	total := 0
	b.Select(b.As("COUNT(*)", "t"))
	q, args := b.Build()
	if err := Db().QueryRow(q, args...).Scan(&total); err != nil {
		return meta, data, err
	}

	// set order by ASC
	b.OrderBy(qo.Order...)

	// set limit
	if qo.Limit != nil {
		b.Limit(lib.Min(*qo.Limit, config.Config().GetInt("param.icpp_max")))
	}

	// set offset
	b.Offset(qo.Offset)

	// set select columns
	b.Select(qo.Fields...)

	// build the sql
	q, args = b.Build()

	// execute query
	rows, err := Db().Query(q, args...)
	if err != nil {
		return meta, data, err
	}
	defer rows.Close()

	// iterate results
	for rows.Next() {

		if err := rows.Scan(s.Addr(&t)...); err != nil {
			return meta, data, err
		}

		// dig... get parent data
		if err := dig(qo); err != nil {
			return meta, data, err
		}

		data = append(data, lib.DeRefPtr(t))

		// Assign t its type's zero value.
		// If not, a dirty value could trespass to
		// to the next result.
		lib.Reset(t)
	}

	// check for iteration errors
	// will be called on deferred rows.Close
	if err := rows.Err(); err != nil {
		return meta, data, err
	}

	// response headers meta
	from := qo.Offset + 1
	to := qo.Offset + len(data)
	meta.Range = fmt.Sprintf("%v-%v/%v", lib.Min(from, total), lib.Min(to, total), total)
	if qo.Checksum == 1 {
		bytes, _ := json.Marshal(data)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = fmt.Sprint(checksum)
	}

	return meta, data, nil
}
