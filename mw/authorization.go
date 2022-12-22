package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

// Verify if user.User is stored in the context
// as key/value pair under the "User" key, meaning the user
// was successfuly authenticated. If so, validates
// if said user has a valid entry in the ACL map for
// the requested endpoint.
func Authorization() gin.HandlerFunc {

	return func(c *gin.Context) {

		if u, ok := c.Get("User"); !ok {

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				msg.Get("5"),
			)

		} else if u, ok := u.(ds.User); !ok {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("5"),
			)

		} else {

			g := ds.Grant{
				Role:   u.Role,
				Route:  c.FullPath(),
				Method: c.Request.Method,
			}

			if valid := g.Valid(); !valid {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("8"),
				)
			}

			c.Next()
		}

	}
}
