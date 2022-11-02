package db

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/msg"
)

func IDbind(c *gin.Context, t Table) (map[string]string, *ParamError) {

	var (
		m = make(map[string]string)
		k = Pk(t)
		v = strings.Split(c.Param("id"), ",")
	)

	if len(k) != len(v) {
		e := ParamError{Message: msg.Get("26")}
		return m, &e
	}

	for i, j := range k {
		m[j] = v[i]
	}
	return m, nil
}
