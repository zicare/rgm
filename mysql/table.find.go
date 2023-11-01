package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
)

// Find returns the qo.DataSource record that matches qo settings.
// Supports BeforeSelect(qo) and AfterSelect(qo). AfterSelect allows parent data retrieval through dig params.
// If a parent resource is not found, Find is aborted with a NotFoundError.
// Beware that qo.DataSource must implement ITable.
func (Table) Find(qo *ds.QueryOptions) (meta ds.ResultSetMeta, data interface{}, err error) {

	t, ok := qo.DataSource.(ITable)
	if !ok {
		return meta, data, new(NotITableError)
	}

	s := sqlbuilder.NewStruct(qo.DataSource)
	b := s.SelectFrom(qo.DataSource.Name())

	// set before select constraints
	if params, err := t.BeforeSelect(qo); err != nil {
		return meta, data, new(ds.NotAllowedError)
	} else {
		for k, v := range params {
			b.Where(b.Equal(k, v))
		}
	}

	// set where Equal for Primary param
	for k, v := range qo.Equal[ds.Primary] {
		b.Where(b.Equal(k, v))
	}

	// set where Equal for Url param
	for k, v := range qo.Equal[ds.Url] {
		b.Where(b.Equal(k, v))
	}

	// build the sql
	q, args := b.Build()

	// execute query
	if err := Db().QueryRow(q, args...).Scan(s.Addr(&t)...); err == sql.ErrNoRows {
		return meta, data, new(ds.NotFoundError)
	}

	// run after select
	if err := t.AfterSelect(qo); err != nil {
		return meta, data, err
	}

	// Response headers meta
	if qo.Checksum == 1 {
		bytes, _ := json.Marshal(qo.DataSource)
		checksum := crc32.ChecksumIEEE([]byte(bytes))
		meta.Checksum = fmt.Sprint(checksum)
	}

	return meta, qo.DataSource, nil
}
