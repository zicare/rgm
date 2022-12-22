package mysql

import (
	"github.com/zicare/rgm/ds"
)

// ITable defines an interface for db table access.
type ITable interface {

	// ITable interfaces also fulfills ds.IDataStore
	ds.IDataStore

	// Must attach foreign table data if available
	Dig(f ...string) []Dig

	// BeforeSelect offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the select by returning a *NotAllowedError.
	BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)

	// BeforeInsert offers a chance to complete extra validations, alter values,
	// or even abort the insert by returning an error.
	// Consider using *NotAllowedError and/or validator.validationErrors as these
	// will be treated as such by ctrl.CrudController, others will be treated
	// as an InternalServerError.
	BeforeInsert(qo *ds.QueryOptions) error

	// BeforeDelete offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the select by returning a *NotAllowedError.
	BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)
}

type Table struct{} // Dig exported

func (Table) Dig(f ...string) []Dig {

	return []Dig{}
}

func (Table) BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}

func (Table) BeforeInsert(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}
