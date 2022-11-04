package db

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

//ResultSetMeta exported
type ResultSetMeta struct {
	Range    string
	Checksum string
}

func Fetch(c *gin.Context, tbl Table) (ResultSetMeta, []interface{}, error) {

	var (
		fo      = GetFetchOptions(c, tbl)
		meta    = ResultSetMeta{Range: "*/*", Checksum: "*"}
		total   int
		results []interface{}
		ms      = sqlbuilder.NewStruct(tbl).For(sqlbuilder.MySQL)
		sb      = ms.SelectFrom(tbl.Name())
	)

	// set where scope
	for k, v := range tbl.Scope(c) {
		sb.Where(sb.Equal(k, v))
	}

	// set where IsNull
	for _, j := range fo.IsNull {
		sb.Where(sb.IsNull(j))
	}

	// set where IsNotNull
	for _, j := range fo.IsNotNull {
		sb.Where(sb.IsNotNull(j))
	}

	// set where In
	for k, v := range fo.In {
		sb.Where(sb.In(k, v...))
	}

	// set where NotIn
	for k, v := range fo.NotIn {
		sb.Where(sb.NotIn(k, v...))
	}

	// set where Equal
	for k, v := range fo.Equal {
		sb.Where(sb.Equal(k, v))
	}

	// set where NotEqual
	for k, v := range fo.NotEqual {
		sb.Where(sb.NotEqual(k, v))
	}

	// set where GreaterThan
	for k, v := range fo.GreaterThan {
		sb.Where(sb.GreaterThan(k, v))
	}

	// set where GreaterEqualThan
	for k, v := range fo.GreaterEqualThan {
		sb.Where(sb.GreaterEqualThan(k, v))
	}

	// set where LessThan
	for k, v := range fo.LessThan {
		sb.Where(sb.LessThan(k, v))
	}

	// set where LessEqualThan
	for k, v := range fo.LessEqualThan {
		sb.Where(sb.LessEqualThan(k, v))
	}

	// get total count
	sb.Select(sb.As("COUNT(*)", "t"))
	sql, args := sb.Build()
	if err := db.QueryRow(sql, args...).Scan(&total); err != nil {
		//Server error: %s
		return meta, results, msg.Get("25").SetArgs(err.Error()).M2E()
	}

	// set order by
	sb.OrderBy(fo.Order...)

	// set limit
	sb.Limit(fo.Limit)

	// set offset
	sb.Offset(fo.Offset)

	// set select columns
	sb.Select(fo.Column...)

	// build the sql
	sql, args = sb.Build()

	// execute query
	rows, err := Db().Query(sql, args...)
	if err != nil {
		return meta, results, msg.Get("25").SetArgs(err).M2E()
	}
	defer rows.Close()

	// iterate results
	for rows.Next() {
		if err := rows.Scan(ms.AddrWithCols(fo.Column, &tbl)...); err != nil {
			//Server error: %s
			return meta, results, msg.Get("25").SetArgs(err).M2E()
		}
		//tbl.Dig(c)
		results = append(results, lib.DeRefPtr(tbl))
	}

	// check for iteration errors
	// will be called on deferred rows.Close
	if err := rows.Err(); err != nil {
		//Server error: %s
		return meta, results, msg.Get("25").SetArgs(err).M2E()
	}

	// Response headers meta
	from := fo.Offset + 1
	to := fo.Offset + len(results)
	meta.Range = fmt.Sprintf("%v-%v/%v", from, to, total)
	if fo.Checksum == 1 {
		bytes, _ := json.Marshal(results)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = strconv.FormatUint(uint64(checksum), 16)
	}

	return meta, results, nil
}
