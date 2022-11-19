package mw

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/tps"
)

// Verify if auth.User is stored in the context
// as key/value pair under the "User" key, meaning the user
// was successfuly authenticated. If so, validates
// if said user TPS is being abused.
func Abuse() gin.HandlerFunc {

	return func(c *gin.Context) {

		if u, ok := c.Get("User"); !ok {
			// Not enough permissions
			c.AbortWithStatusJSON(
				401,
				gin.H{"message": msg.Get("8")},
			)
		} else if u, ok := u.(auth.User); !ok {
			// Not enough permissions
			c.AbortWithStatusJSON(
				401,
				gin.H{"message": msg.Get("8")},
			)
		} else if tps.IsEnabled() {
			// Check for abuse
			if date := tps.Transaction(u.Type, u.UID, u.TPS); date != nil && date.After(time.Now()) {
				// TPS limit exceeded
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": msg.Get("10").SetArgs(date)},
				)
			}
		}

		c.Next()
	}
}
