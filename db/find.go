package db

import (
	"database/sql"
	"encoding/json"
	"hash/crc32"
	"strconv"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/msg"
)

// FindResultSetMeta exported
type FindResultSetMeta struct {
	Checksum string
}

// Find exported
func Find(fo *FindOptions) (FindResultSetMeta, interface{}, error) {

	var (
		meta = FindResultSetMeta{Checksum: "*"}
		ms   = sqlbuilder.NewStruct(fo.Table).For(sqlbuilder.MySQL)
		sb   = ms.SelectFrom(fo.Table.Name())
	)

	// set where scope
	for k, v := range fo.Table.Scope(fo.UID, fo.Parents...) {
		sb.Where(sb.Equal(k, v))
	}

	// set where Equal
	for k, v := range fo.Where {
		sb.Where(sb.Equal(k, v))
	}

	// build the sql
	q, args := sb.Build()

	// execute query
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&fo.Table)...); err == sql.ErrNoRows {
		e := NotFoundError{Message: msg.Get("18")} //Not found!
		return meta, fo.Table, &e
	} else if err != nil {
		//Server error: %s
		return meta, fo.Table, msg.Get("25").SetArgs(err.Error()).M2E()
	}

	if fo.Dig == 1 {
		fo.Table.Dig()
	}

	// Response headers meta
	if fo.Checksum == 1 {
		bytes, _ := json.Marshal(fo.Table)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = strconv.FormatUint(uint64(checksum), 16)
	}

	return meta, fo.Table, nil
}
