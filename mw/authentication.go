package mw

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/msg"
)

// BasicAuthentication executes HTTP basic authentication.
// If passed, a new key/value pair is stored in the request context.
// key: "User"
// value: auth.User
// In order to pass user's system access date range must be valid.
// Client code can implement custom UserDS or use auth.UserTable.
func BasicAuthentication(ds auth.UserDS) gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get usr and pwd from http request headers
		username, password, ok := c.Request.BasicAuth()
		if ok == false {
			// HTTP basic authentication required
			c.AbortWithStatusJSON(
				401,
				gin.H{"message": msg.Get("3")},
			)
		}

		if u, e := ds.GetUser(username, password); e != nil {
			switch e.(type) {
			case *auth.InvalidCredentials:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			case *auth.ExpiredCredentials:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			default:
				c.AbortWithStatusJSON(
					500,
					gin.H{"message": e},
				)
			}
		} else {
			c.Set("User", u)
			c.Next()
		}

	}
}

// JWTAuthentication executes JWT authentication.
// Token must be correct and not expired.
// Authorization towards the ACL must be passed.
// User can't exceed her TPS rate.
// JWT can't be found in acl.RevokedJWTMap registry.
func JWTAuthentication() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Verify JWT authorization header is properly set
		token := strings.Split(c.GetHeader("Authorization"), " ")
		if (len(token) != 2) || (token[0] != "JWT") {
			// JWT authorization header malformed
			c.AbortWithStatusJSON(
				401,
				gin.H{"message": msg.Get("7")},
			)
			return
		}

		if payload, e := jwt.Decode(token[1]); e != nil {
			switch e.(type) {
			case *jwt.InvalidToken:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			case *jwt.InvalidTokenPayload:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			case *jwt.TamperedToken:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			case *jwt.ExpiredToken:
				c.AbortWithStatusJSON(
					401,
					gin.H{"message": e},
				)
			default:
				// Something went wrong
				c.AbortWithStatusJSON(
					500,
					gin.H{"message": e},
				)
			}
		} else if revoked := jwt.IsRevoked(payload); revoked {
			//Token revoked
			c.AbortWithStatusJSON(
				401,
				gin.H{"message": msg.Get("32")},
			)
		} else {
			c.Set("User",
				auth.User{
					UID:  payload.UID,
					Type: payload.Type,
					Role: payload.Role,
					TPS:  payload.TPS,
					From: payload.Iat,
					To:   payload.Exp,
				})
			c.Next()
		}

	}
}
