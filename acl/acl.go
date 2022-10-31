package acl

import (
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

var acl map[Grant]lib.TimeRange

//Grant exported
type Grant struct {
	Role   int64
	Route  string
	Method string
}

// Loads the acl.ACL map in memory.
// acl.ACL maps all the grants to its corresponding validity
// time window. Grants are made of a role and an endpoint.
// The m param is a reference to a database table
// where grants and validity time window are setup.
// The m db.Table fields representing grant elements
// must be properly annotated with the acl tag.
// The request's endpoint and JWT's role can be matched
// to acl.ACL to allow access or not.
// Init also initializes the revokedJWTMap.
func Init(m db.Table) (err error) {

	var (
		g        Grant
		from, to time.Time
		now      = time.Now()
		sb       = sqlbuilder.MySQL.NewSelectBuilder()
	)

	//get m table's acl-related field names
	f, err := lib.TaggedFields(m, "acl", []string{"role", "route", "method", "from", "to"})
	if err != nil {
		return msg.Get("2").M2E() //ACL tags are not properly set
	}

	//fetch all grants
	sb.From(m.View())
	sb.Select(f...)
	sb.Where(
		sb.LessThan(f[3], now),
		sb.GreaterThan(f[4], now),
	)
	sql, args := sb.Build()
	rows, err := db.Db().Query(sql, args...)
	defer rows.Close()

	//scan rows
	acl = make(map[Grant]lib.TimeRange)
	for rows.Next() {
		err := rows.Scan(&g.Role, &g.Route, &g.Method, &from, &to)
		if err != nil {
			//Server error: %s
			return msg.Get("25").SetArgs(err).M2E()
		}
		acl[g] = lib.TimeRange{From: from, To: to}
	}
	err = rows.Err()
	if err != nil {
		//Server error: %s
		return msg.Get("25").SetArgs(err).M2E()
	}

	//Initialize revokedJWTMap
	jwt.Init()

	//log.Println(acl)
	return nil
}

//Returns the acl.ACL map
func ACL() map[Grant]lib.TimeRange {

	return acl
}
