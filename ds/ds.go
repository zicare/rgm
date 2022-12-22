package ds

type ResultSetMeta struct {
	Range    string
	Checksum string
}

// Interface defines an interface for resource data access.
type IDataStore interface {

	// Must return ds name.
	// i.e. For a db store this is the table name
	Name() string

	Find(qo *QueryOptions) (ResultSetMeta, interface{}, error)

	Fetch(qo *QueryOptions) (ResultSetMeta, []interface{}, error)

	Insert(qo *QueryOptions) error

	Delete(qo *QueryOptions) (int64, error)
}

// UserDSFactory makes a IUserDataStore from a generic dst IDataStore.
type UserDSFactory func(dst IDataStore) (IUserDataStore, error)

// AclDSFactory makes a IAclDataStore from a generic dst IDataStore.
type AclDSFactory func(dst IDataStore) (IAclDataStore, error)

// AclDSFactory makes a IPinDataStore from generic p(pin) and u(user) IDataStore's.
type PinDSFactory func(p, u IDataStore) (IPinDataStore, error)
