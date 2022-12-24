package mw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// BasicAuthentication executes HTTP basic authentication.
// If passed, a new key/value pair is stored in the request context.
// key: "User"
// value: ds.User
func BasicAuthentication(dsrc ds.IUserDataSource, crypto lib.ICrypto) gin.HandlerFunc {

	return func(c *gin.Context) {

		if username, password, ok := c.Request.BasicAuth(); !ok {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("3"),
			)

		} else if u, err := dsrc.Get(username); err != nil {

			switch err.(type) {
			case *ds.InvalidCredentials, *ds.ExpiredCredentials:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					msg.Get("4"),
				)
			default:
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
				)
			}

		} else if !crypto.Compare(password, u.Pwd) {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("4"),
			)

		} else {

			c.Set("User", u)

			c.Next()

		}
	}
}

// JWTAuthentication executes JWT authentication.
// Token must be correct, not expired and not revoked.
// If passed, a new key/value pair is stored in the request context.
// key: "User"
// value: ds.User
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
					msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
				)
			}

		} else if revoked := jwt.IsRevoked(payload); revoked {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				msg.Get("32"),
			)

		} else {

			c.Set("User",
				ds.User{
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
