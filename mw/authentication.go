package mw

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/msg"
)

// BasicAuthentication executes HTTP basic authentication.
// If passed, a new key/value pair is stored in the request context.
// key: "User"
// value: auth.User
// In order to pass user's system access date range must be valid.
// Client code can implement custom UserDS or use auth.UserTable.
func BasicAuthentication(t db.Table) gin.HandlerFunc {

	return func(c *gin.Context) {

		if ds, err := auth.TUserDSFactory(t); err != nil {

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				msg.Get("2").SetArgs("User"),
			)

		} else if username, password, ok := c.Request.BasicAuth(); !ok {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("3"),
			)

		} else if u, err := ds.GetUser(username, password); err != nil {

			switch err.(type) {
			case *db.NotFoundError, *auth.InvalidCredentials:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("4"),
				)
			case *auth.ExpiredCredentials:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("6"),
				)
			case *auth.UserTagsError:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("2").SetArgs("User"),
				)
			default:
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					msg.Get("25").SetArgs(err.Error()),
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

		if token := strings.Split(c.GetHeader("Authorization"), " "); (len(token) != 2) || (token[0] != "JWT") {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("7"),
			)

		} else if payload, err := jwt.Decode(token[1]); err != nil {

			switch err.(type) {
			case *jwt.InvalidToken:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("12"),
				)
			case *jwt.InvalidTokenPayload:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("13"),
				)
			case *jwt.TamperedToken:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("14"),
				)
			case *jwt.ExpiredToken:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("15"),
				)
			default:
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					msg.Get("25").SetArgs(err.Error()),
				)
			}

		} else if revoked := jwt.IsRevoked(payload); revoked {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("32"),
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
