package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

// BaseController exported
type BaseController struct{}

// Hierarchical Find.
// Supports endpoints with nested resources in path.
// If a parent resource is not found, the find is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to determine, whether o not, to abort find based on the parent
// resources, which are also made available within Scope method
// by Find.
func (bc BaseController) Find(c *gin.Context, t db.Table, p ...db.Table) {

	if fos, e := bc.getFindOptions(c, t, p...); e != nil {
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
// by Fetch.
func (bc BaseController) Fetch(c *gin.Context, t db.Table, p ...db.Table) {

	if feo, fios, e := bc.getFetchOptions(c, t, p...); e != nil {
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

func (bc BaseController) getFindOptions(c *gin.Context, t db.Table, p ...db.Table) ([]*db.FindOptions, *db.ParamError) {

	var (
		fos    = []*db.FindOptions{}
		uid    = auth.UID(c)
		qparam = make(db.QParams)
		param  = make(db.Params)
	)

	for k, v := range c.Request.URL.Query() {
		qparam[k] = v
	}

	for _, up := range c.Params {
		param[up.Key] = up.Value
	}

	for _, j := range p {
		if fo, e := db.FindOptionsFactory(j, uid, nil, param, true); e != nil {
			return nil, e
		} else {
			fos = append(fos, fo)
		}
	}

	if fo, e := db.FindOptionsFactory(t, uid, qparam, param, true); e != nil {
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

func (bc BaseController) getFetchOptions(c *gin.Context, t db.Table, p ...db.Table) (*db.FetchOptions, []*db.FindOptions, *db.ParamError) {

	var (
		qparam = make(db.QParams)
		param  = make(db.Params)
		uid    = auth.UID(c)
	)

	for k, v := range c.Request.URL.Query() {
		qparam[k] = v
	}

	for _, up := range c.Params {
		param[up.Key] = up.Value
	}

	feo := db.FetchOptionsFactory(t, uid, qparam)
	fios := []*db.FindOptions{}

	for _, j := range p {
		if fio, e := db.FindOptionsFactory(j, uid, nil, param, true); e != nil {
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
