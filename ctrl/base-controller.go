package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

//BaseController exported
type BaseController struct{}

//Get exported
func (ctrl BaseController) Get(c *gin.Context, tbl db.Table) {

	if err := db.Find(c, tbl); err != nil {
		switch e := err.(type) {
		case *db.ParamError:
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": e},
			)
		case *db.NotFoundError:
			c.JSON(
				http.StatusNotFound,
				gin.H{"message": e},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": e},
			)
		}
	} else {
		c.JSON(http.StatusOK, tbl)
	}
}

//Index exported
func (ctrl BaseController) Index(c *gin.Context, tbl db.Table) {

	if meta, data, err := db.Fetch(c, tbl); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": err},
		)
	} else if len(data) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{"message": msg.Get("18")}, //Not found!
		)
	} else {
		c.Header("X-Range", meta.Range)
		c.Header("X-Checksum", meta.Checksum)
		c.JSON(http.StatusOK, data)
	}
}
