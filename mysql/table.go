package mysql

import (
	"github.com/zicare/rgm/ds"
)

// Table offers default implementation for all ds.IDataSource
// methods, except Name().
// You can always overwrite the methods you need to.
// Consider annonymous embedding of Table in your concrete IDataSource.
type Table struct{}

func (Table) BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}

func (Table) AfterSelect(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeInsert(qo *ds.QueryOptions) *ds.ValidationErrors {
	return nil
}

func (Table) BeforeUpdate(qo *ds.QueryOptions) *ds.ValidationErrors {
	return nil
}

func (Table) BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}
