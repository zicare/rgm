package db

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

// FetchResultSetMeta exported
type FetchResultSetMeta struct {
	Range    string
	Checksum string
}

// Supports nested Fetch.
// Supports parent resources data retrieval.
// If a parent resource is not found, Fetch is aborted with a NotFoundError.
// Resources can implement custom logic in the Scope method
// to impose addtional constraints based on qo.UID, which is intended to hold
// the requesting user id.
func Fetch(qo *QueryOptions) (FetchResultSetMeta, []interface{}, error) {

	var (
		meta    = FetchResultSetMeta{Range: "*/*", Checksum: "*"}
		total   int
		results = []interface{}{}
		ms      = sqlbuilder.NewStruct(qo.Table).For(sqlbuilder.MySQL)
		sb      = ms.SelectFrom(qo.Table.Name())
	)

	// set where scope
	for k, v := range qo.Table.Scope(qo.UID, qo.Parents...) {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal for Url param
	for k, v := range qo.Equal[Url] {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal for Query param
	for k, v := range qo.Equal[Query] {
		sb.Where(sb.Equal(k, v))
	}

	// set where IsNull
	for _, j := range qo.IsNull {
		sb.Where(sb.IsNull(j))
	}

	// set where IsNotNull
	for _, j := range qo.IsNotNull {
		sb.Where(sb.IsNotNull(j))
	}

	// set where In
	for k, v := range qo.In {
		sb.Where(sb.In(k, v...))
	}

	// set where NotIn
	for k, v := range qo.NotIn {
		sb.Where(sb.NotIn(k, v...))
	}

	// set where NotEqual
	for k, v := range qo.NotEqual {
		sb.Where(sb.NotEqual(k, v))
	}

	// set where GreaterThan
	for k, v := range qo.GreaterThan {
		sb.Where(sb.GreaterThan(k, v))
	}

	// set where GreaterEqualThan
	for k, v := range qo.GreaterEqualThan {
		sb.Where(sb.GreaterEqualThan(k, v))
	}

	// set where LessThan
	for k, v := range qo.LessThan {
		sb.Where(sb.LessThan(k, v))
	}

	// set where LessEqualThan
	for k, v := range qo.LessEqualThan {
		sb.Where(sb.LessEqualThan(k, v))
	}

	// get total count
	sb.Select(sb.As("COUNT(*)", "t"))
	sql, args := sb.Build()
	if err := db.QueryRow(sql, args...).Scan(&total); err != nil {
		return meta, results, err
	}

	// set order by
	sb.OrderBy(qo.Order...)

	// set limit
	if qo.Limit != nil {
		sb.Limit(lib.Min(*qo.Limit, config.Config().GetInt("param.icpp_max")))
	}

	// set offset
	sb.Offset(qo.Offset)

	// set select columns
	sb.Select(qo.Column...)

	// build the sql
	sql, args = sb.Build()

	// execute query
	rows, err := Db().Query(sql, args...)
	if err != nil {
		return meta, results, err
	}
	defer rows.Close()

	// iterate results
	for rows.Next() {
		if err := rows.Scan(ms.AddrWithCols(qo.Column, &qo.Table)...); err != nil {
			return meta, results, err
		}
		// dig... get parent data
		if err := dig(qo); err != nil {
			return meta, results, err
		}
		results = append(results, lib.DeRefPtr(qo.Table))
	}

	// check for iteration errors
	// will be called on deferred rows.Close
	if err := rows.Err(); err != nil {
		return meta, results, err
	}

	// response headers meta
	from := qo.Offset + 1
	to := qo.Offset + len(results)
	meta.Range = fmt.Sprintf("%v-%v/%v", lib.Min(from, total), lib.Min(to, total), total)
	if qo.Checksum == 1 {
		bytes, _ := json.Marshal(results)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = strconv.FormatUint(uint64(checksum), 16)
	}

	return meta, results, nil
}
