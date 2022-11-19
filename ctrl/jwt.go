package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/jwt"
)

//JwtController exported
type JwtController struct{}

//Get exported
func (ctrl JwtController) Get(c *gin.Context) {

	if u, ok := c.Get("User"); !ok {
		c.AbortWithStatus(500)
		return
	} else if u, ok := u.(auth.User); !ok {
		c.AbortWithStatus(500)
		return
	} else {
		j := jwt.JWTFactory(u.UID, u.Role, u.Type, u.TPS, u.From, u.To)
		c.JSON(
			http.StatusOK,
			gin.H{
				"token": j.ToString(),
			},
		)
	}

}
