package db

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/msg"
)

// ByID exported
func ByID(tbl Table, idv ...string) error {

	var (
		idk = Pk(tbl)
		ms  = sqlbuilder.NewStruct(tbl).For(sqlbuilder.MySQL)
		sb  = ms.SelectFrom(tbl.Name())
	)

	if len(idk) != len(idv) {
		e := ParamError{msg.Get("26")} // Composite key missuse
		return &e
	}

	for i, j := range idk {
		sb.Where(sb.Equal(j, idv[i]))
	}

	q, args := sb.Build()
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&tbl)...); err == sql.ErrNoRows {
		e := NotFoundError{Message: msg.Get("18")} // Not found!
		return &e
	} else if err != nil {
		// Server error: %s
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}
	return nil
}
