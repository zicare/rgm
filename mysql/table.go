package mysql

import (
	"database/sql"

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
	// in a per Table basis, or abort the select by returning an error.
	BeforeSelect(qo *ds.QueryOptions) (ds.Params, error)

	// AfterSelect offers a chance to optionally modify a selected result and attach parent
	// table data in a per Table basis, or abort the select by returning an error.
	AfterSelect(qo *ds.QueryOptions) error

	// BeforeInsert offers a chance to complete extra validations, alter values,
	// or abort the insert by returning an error.
	// Consider using *ds.NotAllowedError and/or *ds.ValidationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeInsert(qo *ds.QueryOptions, tx *sql.Tx) error

	// AfterInsert offers a chance to complete extra validations, alter values,
	// or abort the insert by returning an error.
	// Consider using *ds.NotAllowedError and/or *ds.ValidationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	AfterInsert(qo *ds.QueryOptions, tx *sql.Tx) error

	// BeforeUpdate offers a chance to complete extra validations, alter values,
	// or abort the update by returning an error.
	// Consider using *ds.NotAllowedError and/or *ds.ValidationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	// Consider using tx for any db modification here.
	BeforeUpdate(qo *ds.QueryOptions, tx *sql.Tx) error

	// AfterUpdate offers a chance to complete extra actions
	// or abort the update by returning an error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	// Consider using tx for any db modification here.
	AfterUpdate(qo *ds.QueryOptions, tx *sql.Tx) error

	// BeforeDelete offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the delete by returning an error.
	// Consider using tx for any db modification here.
	BeforeDelete(qo *ds.QueryOptions, tx *sql.Tx) (ds.Params, error)

	// AfterDelete offers a chance to complete extra actions after the delete
	// or even abort the delete by returning an error.
	// Consider using tx for any db modification here.
	AfterDelete(qo *ds.QueryOptions, tx *sql.Tx) error
}

// Table offers default implementation for all ITable and ds.IDataSource
// methods, except Name().
// Consider annonymous embedding of Table in your concrete ITable.
type Table struct{}

func (Table) BeforeSelect(qo *ds.QueryOptions) (ds.Params, error) {
	return nil, nil
}

func (Table) AfterSelect(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeInsert(qo *ds.QueryOptions, tx *sql.Tx) error {
	return nil
}

func (Table) AfterInsert(qo *ds.QueryOptions, tx *sql.Tx) error {
	return nil
}

func (Table) BeforeUpdate(qo *ds.QueryOptions, tx *sql.Tx) error {
	return nil
}

func (Table) AfterUpdate(qo *ds.QueryOptions, tx *sql.Tx) error {
	return nil
}

func (Table) BeforeDelete(qo *ds.QueryOptions, tx *sql.Tx) (ds.Params, error) {
	return nil, nil
}

func (Table) AfterDelete(qo *ds.QueryOptions, tx *sql.Tx) error {
	return nil
}
