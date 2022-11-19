// The auth package is meant to support authentication and authorization middlewares.
// To speed up its work, an ACL map with all grants is stored in-memory.
// The ACL and User data stores must fulfill the AclDS and UserDS interfaces respectively.
// Default TAclDS and TUserDS implementations are included, in case the data is stored in DB.
package auth

import (
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory access control list.
// acl maps grants to validity time window.
// Helps speed up Authorization middleware.
var acl ACL

// ACL exported
type ACL map[Grant]TimeRange

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

// Validates if g Grant exists, also if
// it is valid at the time.
func (g Grant) Valid() bool {

	now := time.Now()
	if r, ok := acl[g]; !ok {
		return false
	} else if now.Before(r.From) || now.After(r.To) {
		return false
	}
	return true
}

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

// Represents an authenticated user.
// Once passed, included authentication middlewares
// mw.BasicAuthentication and mw.JWTAuthentication,
// store the authenticated User in the gin context key/value
// registry under the key "User". So handlers comming afterwards
// have a convenient access to User.
type User struct {
	UID  string    `json:"uid"`
	Usr  string    `json:"usr"`
	Pwd  string    `json:"pwd"`
	Type string    `json:"type"`
	Role string    `json:"role"`
	TPS  float32   `json:"tps"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Returns the authenticated UID or empty
// string if authentication was skipped.
func UID(c *gin.Context) string {

	if u, exists := c.Get("User"); !exists {
		return ""
	} else if u, ok := u.(User); ok {
		return u.UID
	}
	return ""
}
