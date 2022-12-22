package mw

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/tps"
)

// Verify if user.User is stored in the context
// as key/value pair under the "User" key, meaning the user
// was successfuly authenticated. If so, it validates
// if said user TPS is being abused.
func Abuse() gin.HandlerFunc {

	return func(c *gin.Context) {

		if u, ok := c.Get("User"); !ok {

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				msg.Get("5"),
			)

		} else if u, ok := u.(ds.User); !ok {

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				msg.Get("5"),
			)

		} else if tps.IsEnabled() {

			// Check for abuse
			if date := tps.Transaction(u.Type, u.UID, u.TPS); date != nil && date.After(time.Now()) {
				// TPS limit exceeded
				c.AbortWithStatusJSON(
					http.StatusTooManyRequests,
					msg.Get("10").SetArgs(date),
				)
			}

		} else {

			c.Next()

		}

	}
}
