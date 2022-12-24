package ctrl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

// CrudController exported
type CrudController struct{}

// Find exported
func (cc CrudController) Find(c *gin.Context, d ds.IDataSource) {

	if qo, err := ds.QOFactory(c, d); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if meta, data, err := d.Find(qo); err != nil {

		switch err.(type) {
		case *ds.NotFoundError:
			c.JSON(
				http.StatusNotFound,
				msg.Get("18"),
			)
		case *ds.NotAllowedError:
			c.JSON(
				http.StatusUnauthorized,
				msg.Get("11"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else if c.Request.Method == "HEAD" {

		c.Header("X-Range", meta.Range)
		c.Header("X-Checksum", meta.Checksum)

	} else {

		c.Header("X-Checksum", meta.Checksum)
		c.JSON(
			http.StatusOK,
			data,
		)
	}
}

// Fetch exported
func (cc CrudController) Fetch(c *gin.Context, d ds.IDataSource) {

	if qo, err := ds.QOFactory(c, d); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if meta, data, err := d.Fetch(qo); err != nil {

		switch err.(type) {
		case *ds.NotFoundError:
			c.JSON(
				http.StatusNotFound,
				msg.Get("18"),
			)
		case *ds.NotAllowedError:
			c.JSON(
				http.StatusUnauthorized,
				msg.Get("11"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else if c.Request.Method == "HEAD" {

		c.Header("X-Range", meta.Range)
		c.Header("X-Checksum", meta.Checksum)

	} else {

		c.Header("X-Range", meta.Range)
		c.Header("X-Checksum", meta.Checksum)
		c.JSON(http.StatusOK, data)

	}
}

// Post exported
func (cc CrudController) Post(c *gin.Context, dsrc ds.IDataSource) {

	if qo, err := ds.QOFactory(c, dsrc); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if err := c.ShouldBindJSON(dsrc); err != nil {

		c.JSON(
			http.StatusBadRequest,
			msg.ValidationErrors(err),
		)

	} else if err := dsrc.Insert(qo); err != nil {

		switch err.(type) {
		case *ds.NotAllowedError:
			c.JSON(
				http.StatusUnauthorized,
				msg.Get("11"),
			)
		case validator.ValidationErrors: //, *time.ParseError, *json.UnmarshalTypeError
			// Payload didn't pass Table's BeforeInsert validation.
			c.JSON(
				http.StatusBadRequest,
				msg.ValidationErrors(err),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else {

		c.JSON(
			http.StatusCreated,
			dsrc,
		)

	}
}

// Delete exported
func (cc CrudController) Delete(c *gin.Context, dsrc ds.IDataSource) {

	if qo, err := ds.QOFactory(c, dsrc); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if r, err := dsrc.Delete(qo); err != nil {

		switch err.(type) {
		case *ds.NotAllowedError:
			c.JSON(
				http.StatusUnauthorized,
				msg.Get("11"),
			)
		case *ds.ForeignKeyConstraint:
			c.JSON(
				http.StatusConflict,
				msg.Get("40"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else if r == 0 {

		c.JSON(
			http.StatusNotFound,
			msg.Get("18"),
		)

	} else {

		c.JSON(
			http.StatusOK,
			msg.Get("29").SetArgs(r),
		)

	}

}
