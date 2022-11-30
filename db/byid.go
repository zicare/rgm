package db

import (
	"encoding/json"
)

// Find record by primary key.
// In case of composite primary keys, id params must be entered in the same order
// as pk tags are declared in the Table, top to bottom.
func ByID(t Table, id ...string) error {

	var (
		w  = make(UParams)
		pk = Pk(t)
	)

	if len(pk) != len(id) {
		return new(ParamError)
	}

	for k, v := range pk {
		w[v] = id[k]
	}

	if qo := QueryOptionsFactory(t, "", nil, w); !qo.IsPrimary() {
		return new(ParamError)
	} else if _, err := Find(qo); err != nil {
		return err
	} else {
		data, _ := json.Marshal(qo.Table)
		json.Unmarshal(data, &t)
	}

	return nil
}
