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

// Supports nested finds.
// Supports parent resources data retrieval.
// If a parent resource is not found, Find is aborted with a NotFoundError.
// Resources can implement custom logic in the Scope method
// to impose addtional constraints based on qo.UID, which is intended to hold
// the requesting user id.
func Find(qos ...*QueryOptions) (meta FindResultSetMeta, err error) {

	for _, qo := range qos {
		if meta, err = find(qo); err != nil {
			return meta, err
		}
	}
	return meta, nil
}

// Find exported
func find(qo *QueryOptions) (FindResultSetMeta, error) {

	var (
		meta = FindResultSetMeta{Checksum: "*"}
		ms   = sqlbuilder.NewStruct(qo.Table).For(sqlbuilder.MySQL)
		sb   = ms.SelectFrom(qo.Table.Name())
	)

	// set where scope
	for k, v := range qo.Table.Scope(qo.UID, qo.Parents...) {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal for Primary param
	for k, v := range qo.Equal[Primary] {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal for Url param
	for k, v := range qo.Equal[Url] {
		sb.Where(sb.Equal(k, v))
	}

	// build the sql
	q, args := sb.Build()

	// execute query
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&qo.Table)...); err == sql.ErrNoRows {
		return meta, new(NotFoundError)
	} else if err != nil {
		return meta, err
	}

	// dig...  get parent data
	if err := dig(qo); err != nil {
		return meta, err
	}

	// Response headers meta
	if qo.Checksum == 1 {
		bytes, _ := json.Marshal(qo.Table)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = strconv.FormatUint(uint64(checksum), 16)
	}

	return meta, nil
}
