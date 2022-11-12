package ctrl

import (
	"fmt"
	"net/http"

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

// Hierarchical Find.
// Supports endpoints with nested resources in path.
// If a parent resource is not found, the find is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to determine, whether o not, to abort find based on the parent
// resources, which are also made available within Scope method
// by HFind.
func (bc BaseController) HFind(c *gin.Context, t db.Table, p ...db.Table) {

	if fos, e := bc.getHFindOptions(c, t, p...); e != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": e},
		)
	} else {
		for i, fo := range fos {
			if meta, data, err := db.Find(fo); err != nil {
				switch err.(type) {
				case *db.NotFoundError:
					c.JSON(
						http.StatusNotFound,
						gin.H{"message": err},
					)
				default:
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"message": err},
					)
				}
				return
			} else if i == len(fos)-1 {
				c.Header("X-Checksum", meta.Checksum)
				c.JSON(http.StatusOK, data)
			}
		}
	}
}

// Hierarchical Fetch.
// Supports endpoints with nested resources in path.
// If a parent resource is not found, the fetch is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to determine, whether o not, to abort fetch based on the parent
// resources, which are also made available within Scope method
// by HFetch.
func (bc BaseController) HFetch(c *gin.Context, t db.Table, p ...db.Table) {

	if feo, fios, e := bc.getHFetchOptions(c, t, p...); e != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": e},
		)
	} else {
		for _, fio := range fios {
			if _, _, err := db.Find(fio); err != nil {
				switch err.(type) {
				case *db.NotFoundError:
					c.JSON(
						http.StatusNotFound,
						gin.H{"message": err},
					)
				default:
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"message": err},
					)
				}
				return
			}
		}
		if meta, data, err := db.Fetch(feo); err != nil {
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
}

func (bc BaseController) getFindOptions(c *gin.Context, tbl db.Table) (*db.FindOptions, *db.ParamError) {

	var (
		uid    = fmt.Sprint(acl.UserID(c))
		qparam = c.Request.URL.Query()
		uparam = make(map[string]string)
	)

	for _, up := range c.Params {
		uparam[up.Key] = up.Value
	}

	return db.FindOptionsFactory(tbl, uid, qparam, uparam)
}

func (bc BaseController) getFetchOptions(c *gin.Context, tbl db.Table) *db.FetchOptions {

	var (
		uid    = fmt.Sprint(acl.UserID(c))
		qparam = c.Request.URL.Query()
	)

	return db.FetchOptionsFactory(tbl, uid, qparam)
}

func (bc BaseController) getHFindOptions(c *gin.Context, t db.Table, p ...db.Table) ([]*db.FindOptions, *db.ParamError) {

	var (
		fos    = []*db.FindOptions{}
		uid    = fmt.Sprint(acl.UserID(c))
		qparam = c.Request.URL.Query()
		uparam = make(map[string]string)
	)

	for _, up := range c.Params {
		uparam[up.Key] = up.Value
	}

	for _, j := range p {
		if fo, e := db.FindOptionsFactory(j, uid, nil, uparam); e != nil {
			return nil, e
		} else {
			fos = append(fos, fo)
		}
	}

	if fo, e := db.FindOptionsFactory(t, uid, qparam, uparam); e != nil {
		return nil, e
	} else {
		fos = append(fos, fo)
	}

	// append parents
	for i, fo1 := range fos {
		for j, fo2 := range fos {
			if j > i {
				fo2.Parents = append(fo2.Parents, fo1.Table)
			}
		}
	}

	return fos, nil
}

func (bc BaseController) getHFetchOptions(c *gin.Context, t db.Table, p ...db.Table) (*db.FetchOptions, []*db.FindOptions, *db.ParamError) {

	var (
		uid    = fmt.Sprint(acl.UserID(c))
		qparam = c.Request.URL.Query()
		feo    = db.FetchOptionsFactory(t, uid, qparam)
		fios   = []*db.FindOptions{}
		uparam = make(map[string]string)
	)

	for _, up := range c.Params {
		uparam[up.Key] = up.Value
	}

	for _, j := range p {
		if fio, e := db.FindOptionsFactory(j, uid, nil, uparam); e != nil {
			return nil, nil, e
		} else {
			fios = append(fios, fio)
		}
	}

	// append parents
	for i, fio1 := range fios {
		feo.Parents = append(feo.Parents, fio1.Table)
		for j, fio2 := range fios {
			if j > i {
				fio2.Parents = append(fio2.Parents, fio1.Table)
			}
		}
	}

	return feo, fios, nil
}
