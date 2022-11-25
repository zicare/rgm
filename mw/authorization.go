package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/msg"
)

// Verify if auth.User is stored in the context
// as key/value pair under the "User" key, meaning the user
// was successfuly authenticated. If so, validates
// if said user has a valid entry in the ACL map for
// the requested endpoint.
func Authorization() gin.HandlerFunc {

	return func(c *gin.Context) {

		if u, ok := c.Get("User"); !ok {

			//Not enough permissions
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"message": msg.Get("8")},
			)

		} else if u, ok := u.(auth.User); !ok {

			//Not enough permissions
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"message": msg.Get("8")},
			)

		} else {

			g := auth.Grant{
				Role:   u.Role,
				Route:  c.FullPath(),
				Method: c.Request.Method,
			}

			if valid := g.Valid(); !valid {
				//Not enough permissions
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"message": msg.Get("8")},
				)
			}
		}

		c.Next()
	}
}
