package db

import (
	"reflect"

	"github.com/huandu/go-sqlbuilder"
)

//Insert exported
func Insert(qo *QueryOptions) error {

	if err := qo.Table.BeforeInsert(qo.UID, qo.Parents...); err != nil {
		return err
	}

	var (
		t = reflect.ValueOf(qo.Table).Elem()
		f []string
		v []interface{}
	)

	for i := 0; i < t.NumField(); i++ {
		if col, dbOk := t.Type().Field(i).Tag.Lookup("db"); dbOk && (col != "-") {
			if _, viewOk := t.Type().Field(i).Tag.Lookup("view"); !viewOk {
				f = append(f, col)
				v = append(v, t.Field(i).Interface())
			}
		}
	}

	ib := sqlbuilder.MySQL.NewInsertBuilder()
	ib.InsertInto(qo.Table.Name())
	ib.Cols(f...)
	ib.Values(v...)

	sql, args := ib.Build()
	//fmt.Println(sql, args)

	ms := sqlbuilder.NewStruct(qo.Table).For(sqlbuilder.MySQL)
	if err := Db().QueryRow(sql+" RETURNING *", args...).
		Scan(ms.AddrWithCols(f, &qo.Table)...); err != nil {
		return err
	}

	//return find(c, m, PID(m, fields.Primary), false)

	return nil

}
