package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"

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
	} else if err != nil {
		return meta, data, err
	}

	// NEW: hydrate requested digs (parents) into the base record
	if err := hydrateDigsOne(qo); err != nil {
		return meta, data, err
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

// hydrateDigsOne populates qo.DataSource's diggable fields for a single base row.
func hydrateDigsOne(qo *ds.QueryOptions) error {
	if len(qo.Dig) == 0 {
		return nil
	}
	br := reflect.Indirect(reflect.ValueOf(qo.DataSource))

	for _, dig := range qo.Dig {
		// collect FK values from precomputed field indexes
		vals := make([]interface{}, len(dig.FKFieldIdx))
		skip := false
		for i, fi := range dig.FKFieldIdx {
			v := br.Field(fi)
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					skip = true
					break
				}
				v = v.Elem()
			}
			vals[i] = v.Interface()
		}
		if skip {
			continue // no parent → best-effort skip
		}

		// query target by PKs = collected FKs
		s2 := sqlbuilder.NewStruct(dig.Target)
		b2 := s2.SelectFrom(dig.Target.Name())
		b2.Select(s2.Columns()...)
		for i, pk := range dig.PK {
			b2.Where(b2.Equal(pk, vals[i]))
		}
		q2, a2 := b2.Build()

		if err := Db().QueryRow(q2, a2...).Scan(s2.Addr(&dig.Target)...); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return err
		}

		// set the dig field on the base struct
		br.Field(dig.FieldIndex).Set(reflect.ValueOf(dig.Target))
	}
	return nil
}
