package db

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/msg"
)

func Find(c *gin.Context, tbl Table) error {

	var (
		where, err = GetFindOptions(c, tbl)
		ms         = sqlbuilder.NewStruct(tbl).For(sqlbuilder.MySQL)
		sb         = ms.SelectFrom(tbl.Name())
	)

	if err != nil {
		//ParamError
		return err
	}

	for k, v := range tbl.Scope(c) {
		sb.Where(sb.Equal(k, v))
	}

	for k, v := range where {
		sb.Where(sb.Equal(k, v))
	}

	q, args := sb.Build()
	//log.Println(q, args)
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&tbl)...); err == sql.ErrNoRows {
		e := NotFoundError{Message: msg.Get("18")} //Not found!
		return &e
	} else if err != nil {
		//Server error: %s
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}

	//tbl.Dig(c)

	return nil
}
