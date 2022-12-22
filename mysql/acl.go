package mysql

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

// MySQL implementation of acl.IAclDataStore.
type aclDataStore struct {
	t ITable
	f []string
}

// UserDSFactory returns an object that implements user.IUserDataStore.
func AclDSFactory(acl ds.IDataStore) (ds.IAclDataStore, error) {

	dst := aclDataStore{}

	t, ok := acl.(ITable)
	if !ok {
		return dst, new(NotITableError)
	}

	// Verify user tags
	if f, err := ds.TagValuesPivoted(t, "db", "json", []string{"role", "route", "method", "from", "to"}); err != nil {
		err.Copy(msg.Get("2").SetArgs("ACL"))
		return dst, err
	} else {
		dst.f = f
		dst.t = t
	}

	return dst, nil
}

// GetAcl exported
func (dst aclDataStore) Fetch() (ds.Acl, error) {

	m := make(ds.Acl)

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(dst.t.Name())
	sb.Select(dst.f...)
	q, args := sb.Build()

	rows, err := Db().Query(q, args...)
	if err != nil {
		return m, err
	}
	defer rows.Close()

	for rows.Next() {
		g := ds.Grant{}
		t := ds.TimeRange{}
		if err := rows.Scan(&g.Role, &g.Route, &g.Method, &t.From, &t.To); err != nil {
			return m, err
		}
		m[g] = t
	}

	return m, nil
}
