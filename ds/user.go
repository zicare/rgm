package ds

import (
	"time"

	"github.com/gin-gonic/gin"
)

// User exported
type User struct {
	UID  string    `json:"uid"`
	Usr  string    `json:"usr"`
	Pwd  string    `json:"pwd"`
	Type string    `json:"type"`
	Role string    `json:"role"`
	TPS  float32   `json:"tps"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Defines an interface to retrieve user data
// and modify password.
type IUserDataStore interface {

	// Return the active User matching the username
	Get(username string) (User, error)

	// Patch password for matching username
	//PatchPwd(patch Patch) error
}

// Returns the authenticated UID or empty
// string if authentication was skipped.
func UID(c *gin.Context) string {

	if u, exists := c.Get("User"); !exists {
		return ""
	} else if u, ok := u.(User); ok {
		return u.UID
	}
	return ""
}
