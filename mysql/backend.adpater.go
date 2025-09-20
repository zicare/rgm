package mysql

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

// backend implements ds.Backend using the existing mysql.Init and mysql.Db.
type backend struct{}

func (b *backend) Meta() (map[string]ds.TableMeta, error) {

	// Tables + Views, keep column order.
	const colSQL = `
SELECT t.TABLE_NAME, t.TABLE_TYPE, c.COLUMN_NAME
FROM information_schema.TABLES  AS t
JOIN information_schema.COLUMNS AS c
  ON c.TABLE_SCHEMA = t.TABLE_SCHEMA
 AND c.TABLE_NAME   = t.TABLE_NAME
WHERE t.TABLE_SCHEMA = DATABASE()
  AND t.TABLE_TYPE IN ('BASE TABLE','VIEW')
ORDER BY t.TABLE_NAME, c.ORDINAL_POSITION;
`
	colRows, err := Db().Query(colSQL)
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	colsByTable := make(map[string][]string)
	isView := make(map[string]bool)
	for colRows.Next() {
		var tbl, ttype, col string
		if err := colRows.Scan(&tbl, &ttype, &col); err != nil {
			return nil, err
		}
		colsByTable[tbl] = append(colsByTable[tbl], col)
		isView[tbl] = (ttype == "VIEW")
	}
	if err := colRows.Err(); err != nil {
		return nil, err
	}

	// PKs for tables only.
	const pkSQL = `
SELECT k.TABLE_NAME, k.COLUMN_NAME
FROM information_schema.TABLE_CONSTRAINTS tc
JOIN information_schema.KEY_COLUMN_USAGE k
  ON k.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
 AND k.TABLE_SCHEMA   = tc.TABLE_SCHEMA
 AND k.TABLE_NAME     = tc.TABLE_NAME
WHERE tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
  AND tc.TABLE_SCHEMA    = DATABASE()
ORDER BY k.TABLE_NAME, k.ORDINAL_POSITION;
`
	pkRows, err := Db().Query(pkSQL)
	if err != nil {
		return nil, err
	}
	defer pkRows.Close()

	pksByTable := make(map[string][]string)
	for pkRows.Next() {
		var tbl, col string
		if err := pkRows.Scan(&tbl, &col); err != nil {
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
			// View: PKs == nil
			reg[tbl] = ds.TableMeta{Columns: cols, PKs: nil}
			continue
		}
		// Table: ensure non-nil even if no PKs
		pks := pksByTable[tbl]
		if pks == nil {
			pks = []string{} // empty, non-nil
		}
		reg[tbl] = ds.TableMeta{Columns: cols, PKs: pks}
	}
	return reg, nil
}
