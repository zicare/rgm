package postgres

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// PostgreSQL implementation of user.IUserDataSource.
type userDataSource struct {
	t ITable
	f []string
}

// UserDSFactory returns an object that implements user.IUserDataSource.
func UserDSFactory(user ds.IDataSource) (ds.IUserDataSource, error) {

	dsrc := userDataSource{}

	t, ok := user.(ITable)
	if !ok {
		return dsrc, new(NotITableError)
	}

	// Verify user tags
	if f, err := ds.TagValuesPivoted(t, "db", "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"}); err != nil {
		err.Copy(msg.Get("2").SetArgs("User"))
		return dsrc, err
	} else {
		dsrc.f = f
		dsrc.t = t
	}

	return dsrc, nil
}

// Get exported
func (dsrc userDataSource) Get(username string) (ds.User, error) {

	u := ds.User{Type: dsrc.t.Name()}

	b := sqlbuilder.NewSelectBuilder()
	b.From(dsrc.t.Name())
	b.Select(dsrc.f...)
	b.Where(b.Equal(dsrc.f[3], username))
	q, args := b.BuildWithFlavor(sqlbuilder.PostgreSQL)

	// execute query
	if err := Db().QueryRow(q, args...).Scan(&u.UID, &u.Role, &u.TPS, &u.Usr, &u.Pwd, &u.From, &u.To); err == sql.ErrNoRows {
		return u, new(ds.InvalidCredentials)
	} else if err != nil {
		return u, err
	}

	// verify if credential are expired
	now := time.Now()
	if now.Before(u.From) || now.After(u.To) {
		return u, new(ds.ExpiredCredentials)
	}

	return u, nil

}

// PatchPwd exported
func (dsrc userDataSource) patchPwd(patch *ds.Patch, crypto lib.ICrypto) error {

	b := sqlbuilder.NewUpdateBuilder()
	b.Update(dsrc.t.Name())
	b.Set(b.Assign(dsrc.f[4], crypto.Encode(patch.Password)))
	b.Where(b.Equal(dsrc.f[3], patch.Email))
	q, args := b.BuildWithFlavor(sqlbuilder.PostgreSQL)

	if res, err := Db().Exec(q, args...); err != nil {
		return err
	} else if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return new(ds.UpdateError)
	}

	return nil
}
