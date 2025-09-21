package mysql

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
)

// Fetch returns the qo.DataSource records that match qo settings.
// Supports BeforeSelect(qo) and parent data retrieval through dig params.
// If a parent resource is not found, Fetch is aborted with a NotFoundError.
// Beware that qo.DataSource must implement ITable.
func (Table) Fetch(qo *ds.QueryOptions) (meta ds.ResultSetMeta, data []interface{}, err error) {

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

	// set where Like
	for k, v := range qo.Like {
		b.Where(b.Like(k, v))
	}

	// set where NotLike
	for k, v := range qo.NotLike {
		b.Where(b.NotLike(k, v))
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

	// buffer base rows here
	var buf []ds.IDataSource

	// iterate results
	for rows.Next() {

		if err := rows.Scan(s.Addr(&t)...); err != nil {
			return meta, data, err
		}

		//data = append(data, lib.DeRefPtr(t))

		// clone current row into its own pointer and buffer it
		tv := reflect.Indirect(reflect.ValueOf(t))
		clone := reflect.New(tv.Type())
		clone.Elem().Set(tv)
		buf = append(buf, clone.Interface().(ds.IDataSource))

		// reset reusable t for next scan
		lib.Reset(t)
	}

	// check for iteration errors
	// will be called on deferred rows.Close
	if err := rows.Err(); err != nil {
		return meta, data, err
	}

	// NEW: hydrate all requested digs in one batched shot per dig
	if err := hydrateDigsMany(qo, buf); err != nil {
		return meta, data, err
	}

	// now build the final data slice
	for _, row := range buf {
		if it, ok := row.(ITable); ok {
			if err := it.AfterSelect(qo); err != nil {
				return meta, data, err
			}
		}
		data = append(data, lib.DeRefPtr(row))
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

// hydrateDigsMany hydrates all buffered base rows' diggable fields with one DB trip per dig.
// Best-effort: skips rows with NULL/incomplete FK tuples or when parent not found;
// bubbles up real driver/scan errors.
func hydrateDigsMany(qo *ds.QueryOptions, rows []ds.IDataSource) error {

	if len(qo.Dig) == 0 || len(rows) == 0 {
		return nil
	}

	type tupleKey string
	keyOf := func(vals []interface{}) tupleKey {
		const sep = "\x1F"
		var b strings.Builder
		for i, v := range vals {
			if i > 0 {
				b.WriteString(sep)
			}
			b.WriteString(fmt.Sprint(v))
		}
		return tupleKey(b.String())
	}

	for _, dig := range qo.Dig {
		// 1) collect distinct FK tuples
		tupleToRows := make(map[tupleKey][]int)
		tuples := make([][]interface{}, 0, len(rows))
		seen := make(map[tupleKey]struct{})

		for i, r := range rows {
			br := reflect.Indirect(reflect.ValueOf(r))
			vals := make([]interface{}, len(dig.FKFieldIdx))
			incomplete := false
			for j, fi := range dig.FKFieldIdx {
				fv := br.Field(fi)
				if fv.Kind() == reflect.Ptr {
					if fv.IsNil() {
						incomplete = true
						break
					}
					fv = fv.Elem()
				}
				vals[j] = fv.Interface()
			}
			if incomplete {
				continue
			}

			k := keyOf(vals)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				tuples = append(tuples, vals)
			}
			tupleToRows[k] = append(tupleToRows[k], i)
		}
		if len(tuples) == 0 {
			continue
		}

		// 2) build one query for all tuples
		s2 := sqlbuilder.NewStruct(dig.Target)
		b2 := s2.SelectFrom(dig.Target.Name())
		b2.Select(s2.Columns()...)

		if len(dig.PK) == 1 {
			vals := make([]interface{}, 0, len(tuples))
			for _, t := range tuples {
				vals = append(vals, t[0])
			}
			b2.Where(b2.In(dig.PK[0], vals...))
		} else {
			var ors []string
			for _, t := range tuples {
				ands := make([]string, 0, len(dig.PK))
				for j, pk := range dig.PK {
					ands = append(ands, b2.Equal(pk, t[j]))
				}
				ors = append(ors, b2.And(ands...))
			}
			b2.Where(b2.Or(ors...))
		}

		q2, a2 := b2.Build()

		// targets lives outside so we can assign after the scoped defer returns
		targets := make(map[tupleKey]reflect.Value)

		// 3) run the query in a small scope so defer closes per-iteration
		if err := func() error {
			rows2, err := Db().Query(q2, a2...)
			if err != nil {
				return err
			}
			defer rows2.Close()

			rt := reflect.Indirect(reflect.ValueOf(dig.Target)).Type()
			for rows2.Next() {
				tgt := reflect.New(rt).Interface().(ds.IDataSource)
				if err := rows2.Scan(s2.Addr(&tgt)...); err != nil {
					return err
				}

				trv := reflect.Indirect(reflect.ValueOf(tgt))
				pkVals := make([]interface{}, len(dig.PKFieldIdx))
				for j, fi := range dig.PKFieldIdx {
					v := trv.Field(fi)
					if v.Kind() == reflect.Ptr {
						if v.IsNil() {
							continue
						}
						v = v.Elem()
					}
					pkVals[j] = v.Interface()
				}
				targets[keyOf(pkVals)] = reflect.ValueOf(tgt)
			}
			if err := rows2.Err(); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}

		// 4) assign into each base row that requested this tuple
		for k, idxs := range tupleToRows {
			tgt, ok := targets[k]
			if !ok {
				continue
			}
			for _, i := range idxs {
				br := reflect.Indirect(reflect.ValueOf(rows[i]))
				br.Field(dig.FieldIndex).Set(tgt)
			}
		}
	}

	return nil
}
