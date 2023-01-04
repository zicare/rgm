package ctrl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// PinController exported
type PinController struct{}

// Post saves a pin to PinDataSource and sends it back to requesting user by email.
func (ctrl PinController) Post(c *gin.Context, fn ds.PinDSFactory, p ds.IDataSource, u ds.IDataSource) {

	d := &struct {
		Email string `json:"usr" binding:"required,email"`
	}{}

	if dsrc, err := fn(p, u); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if err := c.ShouldBindJSON(d); err != nil {

		c.JSON(
			http.StatusBadRequest,
			msg.ValidationErrors(err),
		)

	} else if p, err := dsrc.Post(d.Email); err != nil {

		switch err.(type) {
		case *ds.InvalidCredentials, *ds.ExpiredCredentials:
			c.JSON(
				http.StatusAccepted,
				msg.Get("33"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else {

		p.Send()

		c.JSON(
			http.StatusAccepted,
			msg.Get("33"),
		)

	}

}

// Patch updates the password in IUserDataSource.
func (ctrl PinController) Patch(c *gin.Context, fn ds.PinDSFactory, p ds.IDataSource, u ds.IDataSource, crypto lib.ICrypto) {

	dsrc, err := fn(p, u)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)
	}

	pr := ds.PatchReceiver()
	if err := c.ShouldBindJSON(pr); err != nil {
		c.JSON(
			http.StatusBadRequest,
			msg.ValidationErrors(err),
		)
		return
	}

	patch := ds.PatchDecoder(pr)
	if err := dsrc.PatchPwd(patch, crypto); err != nil {

		ml := []msg.Message{}
		switch err.(type) {
		case *ds.InvalidCredentials:
			ml = append(ml, msg.Get("4").SetField("usr"))
		case *ds.ExpiredCredentials:
			ml = append(ml, msg.Get("6").SetField("usr"))
		case *ds.InvalidPinError:
			ml = append(ml, msg.Get("36").SetField("pin"))
		case *ds.ExpiredPinError:
			ml = append(ml, msg.Get("39").SetField("pin"))
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
			return
		}
		c.JSON(
			http.StatusBadRequest,
			ml,
		)

	} else {

		c.JSON(
			http.StatusOK,
			msg.Get("38"),
		)
	}

}
