package mysql

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

// MySQL implementation of acl.IAclDataSource.
type aclDataSource struct {
	t ITable
	f []string
}

// UserDSFactory returns an object that implements user.IUserDataSource.
func AclDSFactory(acl ds.IDataSource) (ds.IAclDataSource, error) {

	dsrc := aclDataSource{}

	t, ok := acl.(ITable)
	if !ok {
		return dsrc, new(NotITableError)
	}

	// Verify user tags
	if f, err := ds.TagValuesPivoted(t, "db", "json", []string{"role", "route", "method", "from", "to"}); err != nil {
		err.Copy(msg.Get("2").SetArgs("ACL"))
		return dsrc, err
	} else {
		dsrc.f = f
		dsrc.t = t
	}

	return dsrc, nil
}

// GetAcl exported
func (dsrc aclDataSource) Fetch() (ds.Acl, error) {

	m := make(ds.Acl)

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(dsrc.t.Name())
	sb.Select(dsrc.f...)
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
