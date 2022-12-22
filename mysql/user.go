package mysql

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// MySQL implementation of user.IUserDataStore.
type userDataStore struct {
	t ITable
	f []string
}

// UserDSFactory returns an object that implements user.IUserDataStore.
func UserDSFactory(user ds.IDataStore) (ds.IUserDataStore, error) {

	dst := userDataStore{}

	t, ok := user.(ITable)
	if !ok {
		return dst, new(NotITableError)
	}

	// Verify user tags
	if f, err := ds.TagValuesPivoted(t, "db", "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"}); err != nil {
		err.Copy(msg.Get("2").SetArgs("User"))
		return dst, err
	} else {
		dst.f = f
		dst.t = t
	}

	return dst, nil
}

// Get exported
func (dst userDataStore) Get(username string) (ds.User, error) {

	u := ds.User{Type: dst.t.Name()}

	b := sqlbuilder.NewSelectBuilder()
	b.From(dst.t.Name())
	b.Select(dst.f...)
	b.Where(b.Equal(dst.f[3], username))
	q, args := b.Build()

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
func (dst userDataStore) patchPwd(patch *ds.Patch) error {

	b := sqlbuilder.NewUpdateBuilder()
	b.Update(dst.t.Name())
	b.Set(b.Assign(dst.f[4], lib.Crypto().Encode(patch.Password)))
	b.Where(b.Equal(dst.f[3], patch.Email))
	q, args := b.Build()

	if res, err := Db().Exec(q, args...); err != nil {
		return err
	} else if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return new(UpdateError)
	}

	return nil
}
