package ctrl

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

// PinController exported
type PinController struct{}

// Post exported
// Save PIN to db and sends it back.
// p is the Table for pins, must have proper pin tags.
// u is the Table for users, must have proper user tags.
// Check auth.TPINDSFactory for more information.
func (ctrl PinController) Post(c *gin.Context, p, u db.Table) {

	var pin auth.PIN

	if tpds, err := auth.TPINDSFactory(p, u); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("2").SetArgs("PIN/User"),
		)

	} else if err := c.ShouldBindJSON(p); err != nil {

		c.JSON(
			http.StatusBadRequest,
			p.ValidationErrors(err),
		)

	} else if data, err := json.Marshal(p); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(err),
		)

	} else if err := json.Unmarshal(data, &pin); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(err),
		)

	} else if _, err := tpds.PostPIN(pin.Email); err != nil {

		switch err.(type) {
		case *auth.InvalidCredentials:
			c.JSON(
				http.StatusAccepted,
				msg.Get("33"),
			)
		case *db.InsertError:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("34"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(err),
			)
		}

	} else {

		c.JSON(
			http.StatusAccepted,
			msg.Get("33"),
		)

	}

}

// Patch exported
func (ctrl PinController) Patch(c *gin.Context, t db.Table) {

	c.JSON(http.StatusOK, t)
}
