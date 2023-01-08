package mysql

import (
	"github.com/zicare/rgm/ds"
)

// ITable defines an interface for db table access.
// Consider annonymous embedding of Table in your concrete ITable.
// Table offers default implementation for all ITable and ds.IDataSource
// methods, except Name().
// You can always overwrite the methods you need to.
// Check Table for more information.
type ITable interface {

	// ITable interfaces must fulfills ds.IDataSource
	ds.IDataSource

	// BeforeSelect offers a chance optionally set additional constraints
	// in a per Table basis, or abort the select by returning a *ds.NotAllowedError.
	BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)

	// AfterSelect offers a chance to optionally modify a selected result and attach parent
	// table data in a per Table basis, or abort the select by returning an error.
	AfterSelect(qo *ds.QueryOptions) error

	// BeforeInsert offers a chance to complete extra validations, alter values,
	// or abort the insert by returning an error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeInsert(qo *ds.QueryOptions) error

	// BeforeUpdate offers a chance to complete extra validations, alter values,
	// or abort the update by returning an error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeUpdate(qo *ds.QueryOptions) error

	// BeforeDelete offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the select by returning a *ds.NotAllowedError.
	BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)
}

// Table offers default implementation for all ITable and ds.IDataSource
// methods, except Name().
// Consider annonymous embedding of Table in your concrete ITable.
type Table struct{}

func (Table) BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}

func (Table) AfterSelect(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeInsert(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeUpdate(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}
