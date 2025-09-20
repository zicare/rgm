package ds

import (
	"reflect"
	"strings"
	"sync"
)

// Cache of db-tag -> field-index, keyed by concrete struct type.
// Treat returned maps as read-only.
var dbTagIndexCache sync.Map // key: reflect.Type, value: map[string]int

func dbTagIndex(rt reflect.Type) map[string]int {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if m, ok := dbTagIndexCache.Load(rt); ok {
		return m.(map[string]int)
	}
	m := make(map[string]int, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if tag, ok := f.Tag.Lookup("db"); ok && tag != "" && tag != "-" {
			m[tag] = i
		}
	}
	dbTagIndexCache.Store(rt, m)
	return m
}

// returns field indexes for the given db tags on the given type.
// ok=false if any tag isn't found.
func fieldIdxForTags(rt reflect.Type, tags []string) (idx []int, ok bool) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	m := dbTagIndex(rt)
	idx = make([]int, len(tags))
	for i, name := range tags {
		fi, hit := m[name]
		if !hit {
			return nil, false
		}
		idx[i] = fi
	}
	return idx, true
}

func (qo *QueryOptions) setDig(qpar qparams) {

	// Reset dig plan
	qo.Dig = nil

	// No ?dig=... keys → nothing to do
	digs, ok := qpar["dig"]
	if !ok {
		return
	}

	// Base struct type (of the current DataSource)
	r := reflect.Indirect(reflect.ValueOf(qo.DataSource))
	rt := r.Type()

	for _, key := range digs {

		// Resolve diggable field index by the json key
		idx, ok := qo.DigableFields[key]
		if !ok {
			continue
		} // !ok means unknown dig key; ignore
		sf := rt.Field(idx)

		// Parse fk tag: fk:"a,b,..." → []string
		var fk []string
		if tag := sf.Tag.Get("fk"); tag != "" {
			for _, s := range strings.Split(tag, ",") {
				if s = strings.TrimSpace(s); s != "" {
					fk = append(fk, s)
				}
			}
		}

		// Target must be *T where *T implements IDataSource
		ft := sf.Type
		if ft.Kind() != reflect.Ptr {
			continue
		}
		tgt, ok := reflect.New(ft.Elem()).Interface().(IDataSource)
		if !ok {
			continue
		}

		// Target PKs (struct order). Enforce arity match with FK.
		tk, _, _, _, te := Meta(tgt)
		if te != nil || len(fk) != len(tk) {
			continue
		}

		// Precompute field indexes once (FK on base; PK on target)
		fkIdx, ok := fieldIdxForTags(rt, fk)
		if !ok {
			continue
		}
		pkIdx, ok := fieldIdxForTags(reflect.Indirect(reflect.ValueOf(tgt)).Type(), tk)
		if !ok {
			continue
		}

		// Record the plan
		qo.Dig = append(qo.Dig, DigSpec{
			Key:        key,
			FieldIndex: idx,
			FK:         fk,
			PK:         tk,
			Target:     tgt,
			FKFieldIdx: fkIdx,
			PKFieldIdx: pkIdx,
		})
	}
}
