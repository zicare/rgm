package ds

type ResultSetMeta struct {
	Range    string
	Checksum string
}

// Interface defines an interface for resource data access.
type IDataSource interface {

	// Must return ds name.
	// i.e. For a db store this is the table name
	Name() string

	Count(qo *QueryOptions) (int64, error)

	Find(qo *QueryOptions) (ResultSetMeta, interface{}, error)

	Fetch(qo *QueryOptions) (ResultSetMeta, []interface{}, error)

	Insert(qo *QueryOptions) error

	Update(qo *QueryOptions) (int64, error)

	Delete(qo *QueryOptions) (int64, error)

	// BeforeSelect offers a chance optionally set additional constraints
	// in a per Table basis, or abort the select by returning a *ds.NotAllowedError.
	BeforeSelect(qo *QueryOptions) (Params, *NotAllowedError)

	// AfterSelect offers a chance to optionally modify a selected result and attach parent
	// table data in a per Table basis, or abort the select by returning an error.
	AfterSelect(qo *QueryOptions) error

	// BeforeInsert offers a chance to complete extra validations, alter values,
	// or abort the insert by returning a *ds.ValidationErrors error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeInsert(qo *QueryOptions) *ValidationErrors

	// BeforeUpdate offers a chance to complete extra validations, alter values,
	// or abort the update by returning a *ds.ValidationErrors error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeUpdate(qo *QueryOptions) *ValidationErrors

	// BeforeDelete offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the select by returning a *ds.NotAllowedError.
	BeforeDelete(qo *QueryOptions) (Params, *NotAllowedError)
}

// UserDSFactory makes a IUserDataSource from a generic dsrc IDataSource.
type UserDSFactory func(dsrc IDataSource) (IUserDataSource, error)

// AclDSFactory makes a IAclDataSource from a generic dsrc IDataSource.
type AclDSFactory func(dsrc IDataSource) (IAclDataSource, error)

// AclDSFactory makes a IPinDataSource from generic p(pin) and u(user) IDataSource's.
type PinDSFactory func(p, u IDataSource) (IPinDataSource, error)
