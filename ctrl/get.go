package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/db"
)

//BaseController exported
type BaseController struct{}

//Get exported
func (ctrl BaseController) Get(c *gin.Context, tbl db.Table) {

	if whereMap, paramError := db.IDbind(c, tbl); paramError != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": paramError.Message},
		)
	} else if err := db.Find(c, tbl, whereMap); err != nil {
		switch e := err.(type) {
		case *db.NotFoundError:
			c.JSON(
				http.StatusNotFound,
				gin.H{"message": e.Message},
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
