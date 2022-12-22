package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/msg"
)

//JwtController exported
type JwtController struct{}

//Get exported
func (ctrl JwtController) Get(c *gin.Context) {

	if u, ok := c.Get("User"); !ok {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("5"),
		)

	} else if u, ok := u.(ds.User); !ok {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("5"),
		)

	} else {

		j := jwt.JWTFactory(u.UID, u.Role, u.Type, u.TPS, u.From, u.To)
		c.JSON(
			http.StatusOK,
			gin.H{"header": j.GetHeader(), "payload": j.GetPayload(), "token": j.ToString()},
		)

	}

}
