package auth

import (
	"encoding/json"

	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

// Defines an interface for ACL data access.
type AclDS interface {

	// Returns all grants mapped to its corresponding validity time range.
	GetAcl() (ACL, error)
}

// Default implementation of UserDS interface,
// suitable for User data stored in a DB.
type TAclDS struct {
	t db.Table
}

// Returns a TAclDS type.
//
// t must be annotated with json tags for "role", "route", "method", "from" and "to" fields,
// otherwise an *AclTagsError will be returned when calling GetAcl() on TAclDS.
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
// ds := TAclDSFactory(new(ACL))
//
// Init(ds)
func TAclDSFactory(t db.Table) (TAclDS, *AclTagsError) {

	var ds TAclDS

	// Verify acl tags
	if _, err := db.TaggedFields(t, "json", []string{"role", "route", "method", "from", "to"}); err != nil {
		// ACL tags are not properly set
		e := AclTagsError{Message: msg.Get("2").SetArgs("ACL")}
		return ds, &e
	}

	ds.t = t
	return ds, nil
}

// GetAcl exported
func (ds TAclDS) GetAcl() (ACL, error) {

	acl := make(ACL)

	if _, rows, err := db.Fetch(db.FetchOptionsFactory(ds.t, "", nil)); err != nil {
		return acl, msg.Get("25").SetArgs(err).M2E()
	} else {
		for _, row := range rows {
			data, _ := json.Marshal(row)
			g := Grant{}
			t := TimeRange{}
			json.Unmarshal(data, &g)
			json.Unmarshal(data, &t)
			acl[g] = t
		}
	}
	return acl, nil
}
