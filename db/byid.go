package db

import (
	"encoding/json"

	"github.com/zicare/rgm/msg"
)

// Find record by primary key.
// In case of composite primary keys, id params must be entered in the same order
// as pk tags are declared in the Table, top to bottom.
func ByID(t Table, id ...string) error {

	var (
		w  = make(Params)
		pk = Pk(t)
	)

	if len(pk) != len(id) {
		e := ParamError{msg.Get("26")} // Composite key missuse
		return &e
	}

	for k, v := range pk {
		w[v] = id[k]
	}

	if fo, err := FindOptionsFactory(t, "", nil, w, true); err != nil {
		return err
	} else if _, data, err := Find(fo); err != nil {
		return err
	} else {
		data, _ := json.Marshal(data)
		json.Unmarshal(data, &t)
	}

	return nil
}
