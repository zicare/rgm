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

// Supports endpoints with nested resources in path.
// If a parent resource is not found, Find is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to impose addtional constraints based on the requesting user UID
// and parent resources, both of which are made available within Scope method
// by Find.
func (bc BaseController) Find(c *gin.Context, t db.Table, p ...db.Table) {

	if qo, pqos, e := bc.getQueryOptions(c, true, t, p...); e != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": msg.Get("26")},
		)
	} else {
		pqos = append(pqos, qo)
		for i, qo := range pqos {
			if meta, data, err := db.Find(qo); err != nil {
				switch err.(type) {
				case *db.NotFoundError:
					c.JSON(
						http.StatusNotFound,
						gin.H{"message": msg.Get("18")},
					)
				default:
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"message": msg.Get("25").SetArgs(err.Error())},
					)
				}
				return
			} else if i == len(pqos)-1 {
				c.Header("X-Checksum", meta.Checksum)
				c.JSON(http.StatusOK, data)
			}
		}
	}
}

// Supports endpoints with nested resources in path.
// If a parent resource is not found, the Fetch is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to impose addtional constraints based on the requesting user UID
// and parent resources, both of which are made available within Scope method
// by Fetch.
func (bc BaseController) Fetch(c *gin.Context, t db.Table, p ...db.Table) {

	if qo, pqos, e := bc.getQueryOptions(c, false, t, p...); e != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": msg.Get("26")},
		)
	} else {
		for _, pqo := range pqos {
			if _, _, err := db.Find(pqo); err != nil {
				switch err.(type) {
				case *db.NotFoundError:
					c.JSON(
						http.StatusNotFound,
						gin.H{"message": msg.Get("18")},
					)
				default:
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"message": msg.Get("25").SetArgs(err.Error())},
					)
				}
				return
			}
		}
		if meta, data, err := db.Fetch(qo); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": msg.Get("25").SetArgs(err.Error())},
			)
		} else {
			c.Header("X-Range", meta.Range)
			c.Header("X-Checksum", meta.Checksum)
			c.JSON(http.StatusOK, data)
		}
	}
}

// Supports single and multiple deletes.
// Supports endpoints with nested resources in path.
// If a parent resource is not found, Delete is aborted with a NotFoundError.
// Child resources can implement custom logic in the Scope method
// to determine, whether o not, to abort Delete based on the parent
// resources finding result.
// Child resources can also implement custom logic in the BeforeDelete method
// to impose addtional constraints based on the requesting user UID
// and parent resources, both of which are made available within BeforeDelete method
// by Delete. BeforeDelete can also return a flag to abort Delete altogether.
func (bc BaseController) Delete(c *gin.Context, t db.Table, p ...db.Table) {

	qo, pqos, e := bc.getQueryOptions(c, false, t, p...)

	if e != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": msg.Get("26")},
		)
		return
	}

	// Verify parent resources exist and are within read scope
	for _, pqo := range pqos {
		if _, _, err := db.Find(pqo); err != nil {
			switch err.(type) {
			case *db.NotFoundError:
				c.JSON(
					http.StatusNotFound,
					gin.H{"message": msg.Get("18")},
				)
			default:
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"message": msg.Get("25").SetArgs(err.Error())},
				)
			}
			return
		}
	}

	// Proceed with delete.
	if r, err := db.Delete(qo); err != nil {
		switch err.(type) {
		case *db.NotAllowedError:
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"message": msg.Get("11")},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": msg.Get("25").SetArgs(err.Error())},
			)
		}
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{"message": msg.Get("29").SetArgs(r)},
		)
	}

}

func (bc BaseController) getQueryOptions(c *gin.Context, pk bool, t db.Table,
	p ...db.Table) (qo *db.QueryOptions, pqos []*db.QueryOptions, e *db.ParamError) {

	var (
		uid  = auth.UID(c)
		qpar = make(db.QParams)
		upar = make(db.UParams)
	)

	for k, v := range c.Request.URL.Query() {
		qpar[k] = v
	}

	for _, up := range c.Params {
		upar[up.Key] = up.Value
	}

	// No need to pass parents here, we will append them down below
	qo = db.QueryOptionsFactory(t, uid, qpar, upar)
	if pk && !qo.IsPrimary() {
		return nil, nil, new(db.ParamError)
	}

	for _, j := range p {
		if pqo := db.QueryOptionsFactory(j, uid, nil, upar); !pqo.IsPrimary() {
			return nil, nil, new(db.ParamError)
		} else {
			pqos = append(pqos, pqo)
		}
	}

	// append parents
	for i, pqo1 := range pqos {
		qo.Parents = append(qo.Parents, pqo1.Table)
		for j, pqo2 := range pqos {
			if j > i {
				pqo2.Parents = append(pqo2.Parents, pqo1.Table)
			}
		}
	}

	return qo, pqos, nil
}
