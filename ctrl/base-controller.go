package ctrl

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/acl"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

// BaseController exported
type BaseController struct{}

// Find exported
func (bc BaseController) Find(c *gin.Context, tbl db.Table) {

	if fo, e := bc.getFindOptions(c, tbl); e != nil {
		// ParamError, most probably a composite pk malformed
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": e},
		)
	} else if meta, data, err := db.Find(fo); err != nil {
		switch e := err.(type) {
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
		c.Header("X-Checksum", meta.Checksum)
		c.JSON(http.StatusOK, data)
	}
}

// Fetch exported
func (bc BaseController) Fetch(c *gin.Context, tbl db.Table) {

	fo := bc.getFetchOptions(c, tbl)
	if meta, data, err := db.Fetch(fo); err != nil {
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

// Find exported
func (bc BaseController) FindFetch(c *gin.Context, t1 db.Table, t2 db.Table) {

	if fio, feo, e := bc.getFindFetchOptions(c, t1, t2); e != nil {
		// ParamError, most probably a composite pk malformed
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message1": e},
		)
	} else if _, _, err := db.Find(fio); err != nil {
		switch e := err.(type) {
		case *db.NotFoundError:
			c.JSON(
				http.StatusNotFound,
				gin.H{"message2": e},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message3": e},
			)
		}
	} else if meta, data, err := db.Fetch(feo); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message4": err},
		)
	} else if len(data) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{"message5": msg.Get("18")}, //Not found!
		)
	} else {
		c.Header("X-Range", meta.Range)
		c.Header("X-Checksum", meta.Checksum)
		c.JSON(http.StatusOK, data)
	}
}

func (bc BaseController) getFindOptions(c *gin.Context, tbl db.Table) (*db.FindOptions, *db.ParamError) {

	var (
		uid   = fmt.Sprint(acl.UserID(c))
		param = c.Request.URL.Query()
		idv   = strings.Split(c.Param("id"), ",")
	)

	return db.FindOptionsFactory(tbl, uid, param, idv)
}

func (bc BaseController) getFetchOptions(c *gin.Context, tbl db.Table) *db.FetchOptions {

	var (
		uid   = fmt.Sprint(acl.UserID(c))
		param = c.Request.URL.Query()
	)

	return db.FetchOptionsFactory(tbl, uid, param)
}

func (bc BaseController) getFindFetchOptions(c *gin.Context, t1 db.Table, t2 db.Table) (*db.FindOptions, *db.FetchOptions, *db.ParamError) {

	var (
		uid   = fmt.Sprint(acl.UserID(c))
		param = c.Request.URL.Query()
		idv   = strings.Split(c.Param("id"), ",")
	)

	fio, err := db.FindOptionsFactory(t1, uid, nil, idv)
	feo := db.FetchOptionsFactory(t2, uid, param)

	feo.Parents = append(feo.Parents, fio.Table)

	return fio, feo, err
}
