package ds

import (
	"reflect"
)

func dsMeta(d IDataSource) (k, f, w []string, v []interface{}, e *TagError) {

	r := reflect.ValueOf(d).Elem()
	for i := 0; i < r.NumField(); i++ {
		if db, ok := r.Type().Field(i).Tag.Lookup("db"); ok && db != "-" {
			f = append(f, db)
			if pk, ok := r.Type().Field(i).Tag.Lookup("pk"); ok && pk == "1" {
				k = append(k, db)
			}
			if vw, ok := r.Type().Field(i).Tag.Lookup("view"); !ok && vw == "1" {
				w = append(w, db)
				v = append(v, r.Field(i).Interface())
			}
		}
	}

	if len(k) == 0 {
		return k, f, w, v, new(TagError)
	}

	return k, f, w, v, nil
}

/*
func dsWrite(d IDataSource) (fld []string, val []interface{}, e *TagError) {

	v := reflect.ValueOf(d).Elem()
	for i := 0; i < v.NumField(); i++ {
		if db, ok := v.Type().Field(i).Tag.Lookup("db"); ok && (db != "-") {
			if _, ok := v.Type().Field(i).Tag.Lookup("view"); !ok {
				fld = append(fld, db)
				val = append(val, v.Field(i).Interface())
			}
		}
	}

	if len(fld) == 0 {
		return fld, val, new(TagError)
	}

	return fld, val, nil
}
*/

// Works on structs with two tag sets, let's call them target and pivot.
//
// Returns a slice with the target tag's values, provided the pivot
// tag name and values are matched.
//
// model is the struct to work with.
//
// tagName is the target tag name.
//
// pivotTagName is the pivot tag name to be matched.
//
// pivotTagFields are the pivot tag values to be matched.
//
// Example:
//
// type User struct {
//	 UserID    int64     `db:"user_id"   auth:"id"    json:"user_id"   pk:"1"`
//	 RoleID    *int64    `db:"role_id"   auth:"role"  json:"role_id"`
//	 Email     string    `db:"email"     auth:"usr"   json:"email"`
// }
//
// fields, _ := TagFieldsPivoted(new(User), "auth", []string{"id","role","usr"})
// fields -> []string{"user_id","role_id","email"}
//
func TagValuesPivoted(dsrc IDataSource, targetTagKey string, pivotTagKey string, pivotTagValues []string) ([]string, *TagError) {

	targetTagValues := make([]string, len(pivotTagValues))

	t := reflect.ValueOf(dsrc).Elem()
	for i := 0; i < t.NumField(); i++ {
		if ptv, ok := t.Type().Field(i).Tag.Lookup(pivotTagKey); ok {
			if ttv, ok := t.Type().Field(i).Tag.Lookup(targetTagKey); ok && ttv != "-" {
				for j, v := range pivotTagValues {
					if v == ptv {
						targetTagValues[j] = ttv
					}
				}
			}
		}
	}

	for _, f := range targetTagValues {
		if f == "" {
			return targetTagValues, new(TagError)
		}
	}

	return targetTagValues, nil
}

/*
func TaggedFields(tbl IDataSource, tagName string, tagValues []string) ([]string, *TagError) {

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
			return f, new(TagError)
		}
	}

	return f, nil
}

// Returns a slice of db fields tagged as `pk:"1"`
func Pk(tbl IDataSource) (f []string) {

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
func Cols(tbl IDataSource) (f []string) {

	t := reflect.ValueOf(tbl).Elem()

	for i := 0; i < t.NumField(); i++ {
		if db, ok := t.Type().Field(i).Tag.Lookup("db"); ok && db != "-" {
			f = append(f, db)
		}
	}
	return f
}

func Binding(tbl IDataSource, field string) string {

	//t.Field(i).SetString("x")

	t := reflect.ValueOf(tbl).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		if json, ok := f.Tag.Lookup("json"); ok && json == field {
			if binding, ok := f.Tag.Lookup("binding"); ok {
				return binding
			}
			return ""
		}
	}
	return ""
}
*/
