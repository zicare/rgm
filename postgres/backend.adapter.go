package postgres

import (
	"github.com/zicare/rgm/ds"
)

// executed upon package import on any consuming app
// even before the main() method kicks in
func init() {
	ds.DeferInit(func() error {
		return ds.SetBackend(&backend{})
	})
}

// backend implements ds.Backend using the existing postgres.Init and postgres.Db.
type backend struct{}

func (b *backend) Meta() (map[string]ds.TableMeta, error) {

	// Columns for tables ('r') + views ('v') in current_schema().
	const colSQL = `
SELECT c.relname  AS table_name,
       c.relkind  AS kind,
       a.attname  AS column_name
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
JOIN pg_attribute a ON a.attrelid = c.oid
WHERE n.nspname = current_schema()
  AND c.relkind IN ('r','v')
  AND a.attnum  > 0
  AND NOT a.attisdropped
ORDER BY c.relname, a.attnum;
`
	colRows, err := Db().Query(colSQL)
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	colsByTable := make(map[string][]string)
	isView := make(map[string]bool)
	for colRows.Next() {
		var tbl, kind, col string
		if err := colRows.Scan(&tbl, &kind, &col); err != nil {
			return nil, err
		}
		colsByTable[tbl] = append(colsByTable[tbl], col)
		isView[tbl] = (kind == "v")
	}
	if err := colRows.Err(); err != nil {
		return nil, err
	}

	// PKs for tables only.
	const pkSQL = `
SELECT c.relname AS table_name,
       a.attname AS column_name,
       x.ord     AS ordinal
FROM pg_index i
JOIN pg_class c     ON c.oid = i.indrelid
JOIN pg_namespace n ON n.oid = c.relnamespace
JOIN LATERAL unnest(i.indkey) WITH ORDINALITY AS x(attnum, ord) ON TRUE
JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = x.attnum
WHERE i.indisprimary
  AND n.nspname = current_schema()
  AND c.relkind = 'r'
ORDER BY c.relname, x.ord;
`
	pkRows, err := Db().Query(pkSQL)
	if err != nil {
		return nil, err
	}
	defer pkRows.Close()

	pksByTable := make(map[string][]string)
	for pkRows.Next() {
		var tbl, col string
		var _ord int
		if err := pkRows.Scan(&tbl, &col, &_ord); err != nil {
			return nil, err
		}
		pksByTable[tbl] = append(pksByTable[tbl], col)
	}
	if err := pkRows.Err(); err != nil {
		return nil, err
	}

	reg := make(map[string]ds.TableMeta, len(colsByTable))
	for tbl, cols := range colsByTable {
		if isView[tbl] {
			reg[tbl] = ds.TableMeta{Columns: cols, PKs: nil} // view
			continue
		}
		pks := pksByTable[tbl]
		if pks == nil {
			pks = []string{} // table w/ no PK
		}
		reg[tbl] = ds.TableMeta{Columns: cols, PKs: pks}
	}
	return reg, nil
}
