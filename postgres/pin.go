package postgres

import (
	"database/sql"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// PostgreSQL implementation of pin.IPinDataSource.
type pinDataSource struct {
	t ITable
	f []string
	u ds.IUserDataSource
}

// PinDSFactory returns an object that implements pin.IPinDataSource.
func PinDSFactory(pin, user ds.IDataSource) (ds.IPinDataSource, error) {

	pdsrc := pinDataSource{}

	t, ok := pin.(ITable)
	if !ok {
		return pdsrc, new(NotITableError)
	} else if udsrc, err := UserDSFactory(user); err != nil {
		return pdsrc, err
	} else {
		pdsrc.u = udsrc
	}

	// Get pin fields
	if f, err := ds.TagValuesPivoted(t, "db", "json", []string{"email", "code", "created", "expiration"}); err != nil {
		err.Copy(msg.Get("2").SetArgs("Pin"))
		return pdsrc, err
	} else {
		pdsrc.f = f
		pdsrc.t = t
	}

	return pdsrc, nil
}

// Post saves a new pin to p.t.
// email param must match an active user record in p.u.
func (p pinDataSource) Post(email string) (ps ds.Pin, err error) {

	// validate email
	if _, err := p.u.Get(email); err != nil {
		return ps, err
	}

	now := time.Now()
	ps = ds.Pin{
		Email:      email,
		Code:       strings.ToUpper(lib.RandString(config.Config().GetInt("account.pins_length"))),
		Created:    now,
		Expiration: now.Add(30 * time.Minute),
	}

	// insert pin
	b := sqlbuilder.NewInsertBuilder()
	b.InsertInto(p.t.Name())
	b.Cols(p.f...)
	b.Values(ps.Email, ps.Code, ps.Created, ps.Expiration)
	q, args := b.BuildWithFlavor(sqlbuilder.PostgreSQL)
	if res, err := Db().Exec(q, args...); err != nil {
		return ps, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return ps, new(ds.InsertError)
	} else if rows != 1 {
		return ps, new(ds.InsertError)
	}

	return ps, nil

}

// PatchPwd updates password in p.u.
// patch.Email must match an active user record in p.u.
// patch.Email, patch.Pin must match an active pin record in p.
func (p pinDataSource) PatchPwd(patch *ds.Patch, crypto lib.ICrypto) error {

	if _, err := p.u.Get(patch.Email); err != nil {
		// *user.InvalidCredentials, *user.ExpiredCredentials
		return err
	} else if _, err := p.get(patch.Email, patch.Pin); err != nil {
		// *pin.InvalidPinError, *pin.ExpiredPinError
		return err
	} else if err := p.u.(userDataSource).patchPwd(patch, crypto); err != nil {
		return err
	}
	return nil
}

func (p pinDataSource) get(email, code string) (ds.Pin, error) {

	ps := ds.Pin{}

	b := sqlbuilder.NewSelectBuilder()
	b.From(p.t.Name())
	b.Select(p.f...)
	b.Where(b.Equal(p.f[0], email), b.Equal(p.f[1], code))
	q, args := b.BuildWithFlavor(sqlbuilder.PostgreSQL)

	// execute query
	if err := Db().QueryRow(q, args...).Scan(&ps.Email, &ps.Code, &ps.Created, &ps.Expiration); err == sql.ErrNoRows {
		return ps, new(ds.InvalidPinError)
	} else if err != nil {
		return ps, err
	}

	now := time.Now()
	if now.Before(ps.Created) || now.After(ps.Expiration) {
		return ps, new(ds.ExpiredPinError)
	}

	return ps, nil
}
