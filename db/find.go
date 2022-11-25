package db

import (
	"database/sql"
	"encoding/json"
	"hash/crc32"
	"strconv"

	"github.com/huandu/go-sqlbuilder"
)

// FindResultSetMeta exported
type FindResultSetMeta struct {
	Checksum string
}

// Find exported
func Find(qo *QueryOptions) (FindResultSetMeta, interface{}, error) {

	var (
		meta = FindResultSetMeta{Checksum: "*"}
		ms   = sqlbuilder.NewStruct(qo.Table).For(sqlbuilder.MySQL)
		sb   = ms.SelectFrom(qo.Table.Name())
	)

	// set where scope
	for k, v := range qo.Table.Scope(qo.UID, qo.Parents...) {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal
	if qo.IsPrimary() {
		for k, v := range qo.Equal[Primary] {
			sb.Where(sb.Equal(k, v))
		}
	} else {
		for k, v := range qo.Equal[Url] {
			sb.Where(sb.Equal(k, v))
		}
	}

	// build the sql
	q, args := sb.Build()

	// execute query
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&qo.Table)...); err == sql.ErrNoRows {
		return meta, qo.Table, new(NotFoundError)
	} else if err != nil {
		return meta, qo.Table, err
	}

	if qo.Dig == 1 {
		qo.Table.Dig()
	}

	// Response headers meta
	if qo.Checksum == 1 {
		bytes, _ := json.Marshal(qo.Table)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = strconv.FormatUint(uint64(checksum), 16)
	}

	return meta, qo.Table, nil
}
