package auth

import (
	"database/sql"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/mail"
)

type PIN struct {
	Email      string    `json:"email"`
	Code       string    `json:"code"`
	Created    time.Time `json:"created"`
	Expiration time.Time `json:"expiration"`
}

func (p PIN) Send() {

	msg := new(mail.Message)
	msg.To = p.Email
	msg.Subject = "Has recibido un PIN"
	msg.Tpl = "pin.tpl"
	msg.Data = struct{ PIN string }{PIN: p.Code}
	msg.Send(1)
}

// PINDS implementations can be used
// to authenticate users by username and password.
type PINDS interface {

	// Post PIN
	PostPIN(email, code string) (PIN, error)

	// Get PIN matching the email and code
	GetPIN(email, code string) (PIN, error)
}

// Default implementation of PINDS interface.
type tPINDS struct {
	p  db.Table
	u  db.Table
	fp []string
	fu []string
}

// Returns a tPINDS type "object".
//
// p is the Table for pins.
// p must be annotated with json tags for "email", "code", "created" and "expiration" fields,
// otherwise an *PINTagsError will be returned when calling PostPIN() on the returned tPINDS "object".
//
// u is the Table for users.
// u must be annotated with json tags for "uid", "role", "tps", "usr", "pwd", "from" and "to" fields,
// otherwise an *UserTagsError will be returned when calling PostPIN() on the returned tPINDS "object".
//
// Example of p db.Table:
//
// type PIN struct {
// 	 db.BaseTable
// 	 PINID      *int64    `db:"pin_id"       json:"pid"        pk:"1"`
// 	 Email      string    `db:"email"        json:"email"`
// 	 PIN        string    `db:"pin"          json:"code"`
// 	 Created    time.Time `db:"created"      json:"created"`
// 	 Expiration time.Time `db:"expiration"   json:"expiration"`
// }
//
// func (PIN) Name() string {
//	 return "pins"
// }
//
// Example of u db.Table:
//
// type Account struct {
// 	db.BaseTable
// 	UserID      string    `db:"user_id"        json:"uid"`
// 	RoleID      string    `db:"role_id"        json:"role"`
// 	TPS         float32   `db:"tps"            json:"tps"`
// 	Email       string    `db:"email"          json:"usr"`
// 	Password    string    `db:"password"       json:"pwd"`
// 	AccessStart time.Time `db:"access_start"   json:"from"`
// 	AccessEnd   time.Time `db:"access_end"     json:"to"`
// }
//
// func (Account) Name() string {
//	 return "accounts"
// }
//
// ds := TPINDSFactory(new(PIN), new(Account))
//
// pin, err := ds.PostPIN("me@email.com")
//
func TPINDSFactory(p db.Table, u db.Table) (tPINDS, error) {

	var ds = tPINDS{}

	ds.p = p
	ds.u = u

	// Get pin fields
	if fp, err := db.TaggedFields(ds.p, "json", []string{"email", "code", "created", "expiration"}); err != nil {
		return ds, new(PINTagsError)
	} else {
		ds.fp = fp
	}

	// Get user fields
	if fu, err := db.TaggedFields(ds.u, "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"}); err != nil {
		return ds, new(UserTagsError)
	} else {
		ds.fu = fu
	}

	return ds, nil
}

// PostPIN exported
func (ds tPINDS) PostPIN(email string) (PIN, error) {

	now := time.Now()
	p := PIN{
		Email:      email,
		Code:       strings.ToUpper(lib.RandString(config.Config().GetInt("account.pins_length"))),
		Created:    now,
		Expiration: now.Add(30 * time.Minute),
	}

	// validate email
	if err := ds.validateEmail(p.Email); err != nil {
		return p, err
	}

	// insert pin
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(ds.p.Name())
	ib.Cols(ds.fp...)
	ib.Values(p.Email, p.Code, p.Created, p.Expiration)
	q, args := ib.Build()
	if res, err := db.Db().Exec(q, args...); err != nil {
		return p, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return p, err
	} else if rows != 1 {
		return p, new(db.InsertError)
	}

	// send email with pin
	p.Send()

	return p, nil

}

// GetPIN exported
func (ds tPINDS) GetPIN(email, code string) (PIN, error) {

	p := PIN{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(ds.p.Name())
	sb.Select(ds.fp...)
	sb.Equal(ds.fp[1], email)
	sb.Equal(ds.fp[2], code)
	q, args := sb.Build()

	// execute query
	if err := db.Db().QueryRow(q, args...).Scan(&p.Email, &p.Code, &p.Created, &p.Expiration); err == sql.ErrNoRows {
		return p, new(InvalidPIN)
	} else if err != nil {
		return p, err
	}

	now := time.Now()
	if now.Before(p.Created) || now.After(p.Expiration) {
		return p, new(ExpiredPIN)
	}

	return p, nil

}

// GetUser exported
func (ds tPINDS) validateEmail(email string) error {

	u := User{Type: ds.u.Name()}

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(ds.u.Name())
	sb.Select(ds.fu...)
	sb.Where(sb.Equal(ds.fu[3], email))
	q, args := sb.Build()

	// execute query
	if err := db.Db().QueryRow(q, args...).Scan(&u.UID, &u.Role, &u.TPS, &u.Usr, &u.Pwd, &u.From, &u.To); err == sql.ErrNoRows {
		return new(InvalidCredentials)
	} else if err != nil {
		return err
	}

	now := time.Now()
	if now.Before(u.From) || now.After(u.To) {
		return new(ExpiredCredentials)
	}

	return nil

}
