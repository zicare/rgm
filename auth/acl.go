package auth

import (
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/db"
)

// ACL exported
type ACL map[Grant]TimeRange

// In-memory access control list.
// acl maps each grant to a time range.
// Helps speed up Authorization middleware.
var acl ACL

// Meant to be executed on startup, Init loads the acl map in memory.
// ds must implement the AclDS interface. If acl info is
// stored in the DB, consider using the default implementation
// TAclDS and TAclDSFactory for a ready-to-use ds.
func Init(ds AclDS) (err error) {

	if acl, err = ds.GetAcl(); err != nil {
		return err
	}

	return nil
}

// Grant exported
type Grant struct {
	Role   string `json:"role"`
	Route  string `json:"route"`
	Method string `json:"method"`
}

//TimeRange exported
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Validates if g Grant exists and is valid at the time.
func (g Grant) Valid() bool {

	now := time.Now()
	if r, ok := acl[g]; !ok {
		return false
	} else if now.Before(r.From) || now.After(r.To) {
		return false
	}
	return true
}

// Defines an interface for ACL data access.
type AclDS interface {

	// Returns all grants mapped to its corresponding validity time range.
	GetAcl() (ACL, error)
}

// Default implementation of UserDS interface,
// suitable for User data stored in a DB.
type tAclDS struct {
	t db.Table
	f []string
}

// Returns a tAclDS type "object".
//
// t must be annotated with json tags for "role", "route", "method", "from" and "to" fields,
// otherwise an *AclTagsError will be returned.
//
// Example of t db.Table:
//
// type Grant struct {
//	 db.BaseTable
//	 RoleID string    `db:"role_id"  json:"role"`
//	 Route  string    `db:"route"    json:"route"`
//	 Method string    `db:"method"   json:"method"`
//	 Start  time.Time `db:"start"    json:"from"`
//	 End    time.Time `db:"end"      json:"to"`
// }
//
// func (Grant) Name() string {
//	 return "view_grants"
// }
//
// aclDS, err := TAclDSFactory(new(Grant))
//
// acl, err := ds.GetAcl()
//
func TAclDSFactory(t db.Table) (tAclDS, *AclTagsError) {

	var ds tAclDS

	ds.t = t

	// Verify acl tags
	if f, err := db.TaggedFields(ds.t, "json", []string{"role", "route", "method", "from", "to"}); err != nil {
		return ds, new(AclTagsError)
	} else {
		ds.f = f
	}

	return ds, nil
}

// GetAcl exported
func (ds tAclDS) GetAcl() (ACL, error) {

	acl := make(ACL)

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(ds.t.Name())
	sb.Select(ds.f...)
	q, args := sb.Build()

	rows, err := db.Db().Query(q, args...)
	if err != nil {
		return acl, err
	}
	defer rows.Close()

	for rows.Next() {
		g := Grant{}
		t := TimeRange{}
		if err := rows.Scan(&g.Role, &g.Route, &g.Method, &t.From, &t.To); err != nil {
			return acl, err
		}
		acl[g] = t
	}

	return acl, nil
}
