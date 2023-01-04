package validation

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/zicare/rgm/ds"
)

var unik validator.Func = func(fl validator.FieldLevel) bool {

	eq := make(ds.Params)
	noteq := make(ds.Params)

	t := reflect.Indirect(reflect.ValueOf(fl.Parent().Interface()))
	for i := 0; i < t.NumField(); i++ {
		if db, ok := t.Type().Field(i).Tag.Lookup("db"); ok {
			if fv, _, ok := fl.GetStructFieldOKAdvanced(t, t.Type().Field(i).Name); ok {
				if pk, _ := t.Type().Field(i).Tag.Lookup("pk"); pk == "1" {
					noteq[db] = fmt.Sprint(fv.Interface())
				} else if t.Type().Field(i).Name == fl.StructFieldName() {
					eq[db] = fmt.Sprint(fv.Interface())
				}
			}
		}
	}

	if dsrc, ok := fl.Parent().Interface().(ds.IDataSource); ok {
		qo := new(ds.QueryOptions)
		qo.DataSource = dsrc
		qo.Equal = make(map[ds.ParamType]ds.Params)
		qo.Equal[ds.Qry] = eq
		qo.NotEqual = noteq
		if count, err := dsrc.Count(qo); err != nil || count != 0 {
			return false
		}
		return true
	}

	return false
}
