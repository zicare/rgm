package mysql

import (
	"github.com/zicare/rgm/ds"
)

// Dig exported
type Dig struct {
	Table   ITable
	Primary ds.Params
}

func NewDig(t ITable, p ds.Params) Dig {

	return Dig{Table: t, Primary: p}
}

func dig(qo *ds.QueryOptions) error {

	t, ok := qo.DataStore.(ITable)
	if !ok {
		return new(NotITableError)
	}

	for _, dig := range t.Dig(qo.Dig...) {
		dqo := qo.Copy(dig.Table, dig.Primary)
		if _, _, err := t.Find(dqo); err != nil {
			return err
		}
	}
	return nil
}
