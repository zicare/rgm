// The auth package is meant to support authentication and authorization middlewares.
// To speed up its work, an ACL map with all grants is stored in-memory.
// The ACL and User data stores must fulfill the AclDS and UserDS interfaces respectively.
// Default TAclDS and TUserDS implementations are included, in case the data is stored in DB.
package auth

import (
	"github.com/gin-gonic/gin"
)

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
