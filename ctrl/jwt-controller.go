package ctrl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/jwt"
)

//JwtController exported
type JwtController struct{}

//Get exported
func (ctrl JwtController) Get(c *gin.Context) {

	var (
		u           jwt.User
		ok          bool
		now         = time.Now()
		secret      = config.Config().GetString("hmac_key")
		duration, _ = time.ParseDuration(config.Config().GetString("jwt_duration"))
	)

	if user, exists := c.Get("User"); !exists {
		c.AbortWithStatus(500)
		return
	} else if u, ok = user.(jwt.User); !ok {
		c.AbortWithStatus(500)
		return
	}

	//adjust token duration for users with credentials to expire before
	//the default token duration
	if u.To.Before(now.Add(duration)) {
		duration = u.To.Sub(now)
	}

	token, expiration := jwt.Token(u, duration, secret)
	c.JSON(http.StatusOK, gin.H{"token": token, "expiration": expiration})
}
