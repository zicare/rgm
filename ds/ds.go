package ds

type ResultSetMeta struct {
	Range    string
	Checksum string
}

// IDataSource defines an interface for resource data access.
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
}

// UserDSFactory makes a IUserDataSource from a generic dsrc IDataSource.
type UserDSFactory func(dsrc IDataSource) (IUserDataSource, error)

// AclDSFactory makes a IAclDataSource from a generic dsrc IDataSource.
type AclDSFactory func(dsrc IDataSource) (IAclDataSource, error)

// AclDSFactory makes a IPinDataSource from generic p(pin) and u(user) IDataSource's.
type PinDSFactory func(p, u IDataSource) (IPinDataSource, error)
