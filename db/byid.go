package db

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/msg"
)

// ByID exported
func ByID(t Table, id ...string) error {

	var (
		pk = Pk(t)
		ms = sqlbuilder.NewStruct(t).For(sqlbuilder.MySQL)
		sb = ms.SelectFrom(t.Name())
	)

	if len(pk) != len(id) {
		e := ParamError{msg.Get("26")} // Composite key missuse
		return &e
	}

	for i, j := range pk {
		sb.Where(sb.Equal(j, id[i]))
	}

	q, args := sb.Build()
	if err := Db().QueryRow(q, args...).Scan(ms.Addr(&t)...); err == sql.ErrNoRows {
		e := NotFoundError{Message: msg.Get("18")} // Not found!
		return &e
	} else if err != nil {
		// Server error: %s
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}
	return nil
}
