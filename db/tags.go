package db

import (
	"reflect"
)

// Returns the TableMeta object for tbl param.
// The underlying tbl struct must be properly annotated
// with Table related tags. TableMeta supports the
// tags: db, pk, view, serial.
func GetTableMeta(tbl Table) (meta TableMeta) {

	t := reflect.ValueOf(tbl).Elem()

	for i := 0; i < t.NumField(); i++ {
		if col, ok := t.Type().Field(i).Tag.Lookup("db"); ok && (col != "-") {
			meta.Fields = append(meta.Fields, col)
			//check for primary
			if pk, _ := t.Type().Field(i).Tag.Lookup("pk"); pk == "1" {
				meta.Primary = append(meta.Primary, col)
			}
			//check for serial
			if serial, _ := t.Type().Field(i).Tag.Lookup("serial"); serial == "1" {
				meta.Serial = append(meta.Serial, col)
			}
			//check for view or writable
			if view, ok := t.Type().Field(i).Tag.Lookup("view"); !ok {
				meta.Writable = append(meta.Writable, col)
			} else if view == "1" {
				meta.View = append(meta.View, col)
			}
		}
	}

	return
}

// TaggedFields returns a slice of db fields sharing a tag name (tn)
// such fields must be tagged as tn:"tv"
// the tagName value is passed in the tn parameter
// tagValues options are passed in the tv parameter
// ie.:
// type User struct {
//	 UserID    int64     `db:"user_id"   auth:"id"    json:"user_id"   pk:"1"`
//	 RoleID    *int64    `db:"role_id"   auth:"role"  json:"role_id"`
//	 Email     string    `db:"email"     auth:"usr"   json:"email"`
// }
// v ar m User
// f, _ := TaggedFields(m, "auth", []string{"id","role","usr"})
// f -> []string{"user_id","role_id","email"}
func TaggedFields(tbl Table, tagName string, tagValues []string) ([]string, *TableTagError) {

	var (
		f = make([]string, len(tagValues))
		t = reflect.ValueOf(tbl).Elem()
	)

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Type().Field(i).Tag.Lookup(tagName); ok {
			if col, ok := t.Type().Field(i).Tag.Lookup("db"); ok {
				for i, v := range tagValues {
					if v == tag {
						f[i] = col
					}
				}
			}
		}
	}

	for _, col := range f {
		if col == "" {
			return f, new(TableTagError)
		}
	}

	return f, nil
}

// Returns a slice of db fields tagged as `pk:"1"`
func Pk(tbl Table) (f []string) {

	t := reflect.ValueOf(tbl).Elem()

	for i := 0; i < t.NumField(); i++ {
		if pk, ok := t.Type().Field(i).Tag.Lookup("pk"); ok && pk == "1" {
			if db, ok := t.Type().Field(i).Tag.Lookup("db"); ok && db != "-" {
				f = append(f, db)
			}
		}
	}

	return f
}

// Returns a slice of db fields tagged as `pk:"1"`
func Cols(tbl Table) (f []string) {

	t := reflect.ValueOf(tbl).Elem()

	for i := 0; i < t.NumField(); i++ {
		if db, ok := t.Type().Field(i).Tag.Lookup("db"); ok && db != "-" {
			f = append(f, db)
		}
	}

	return f
}
