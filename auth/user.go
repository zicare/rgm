package auth

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"golang.org/x/crypto/bcrypt"
)

// Represents an authenticated user.
// Once passed, included authentication middlewares
// mw.BasicAuthentication and mw.JWTAuthentication,
// store the authenticated User in the gin context key/value
// registry under the key "User". So handlers comming afterwards
// have a convenient access to User.
type User struct {
	UID  string    `json:"uid"`
	Usr  string    `json:"usr"`
	Pwd  string    `json:"pwd"`
	Type string    `json:"type"`
	Role string    `json:"role"`
	TPS  float32   `json:"tps"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// UserDS implementations can be used
// to authenticate users by username and password.
type UserDS interface {

	// Return User matching the username and password
	GetUser(username, password string) (User, error)
}

// Default implementation of UserDS interface.
type tUserDS struct {
	t db.Table
	f []string
}

// Returns a tUserDS type "object".
//
// t must be annotated with json tags for "uid", "role", "tps", "usr", "pwd", "from" and "to" fields,
// otherwise an *UserTagsError will be returned when calling GetUser() on the returned tUserDS "object".
//
// Example of t db.Table:
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
// ds := TUserDSFactory(new(Account))
//
// user, err := ds.GetUser("admin", "secret")
//
func TUserDSFactory(t db.Table) (tUserDS, *UserTagsError) {

	var ds tUserDS

	ds.t = t

	// Verify user tags
	if f, err := db.TaggedFields(ds.t, "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"}); err != nil {
		return ds, new(UserTagsError)
	} else {
		ds.f = f
	}

	return ds, nil
}

// GetUser exported
func (ds tUserDS) GetUser(username, password string) (User, error) {

	u := User{Type: ds.t.Name()}

	// Verify user tags
	f, err := db.TaggedFields(ds.t, "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"})
	if err != nil {
		return u, new(UserTagsError)
	}

	sb := sqlbuilder.NewSelectBuilder()
	sb.From(ds.t.Name())
	sb.Select(f...)
	sb.Where(sb.Equal(f[3], username))
	q, args := sb.Build()

	// execute query
	if err := db.Db().QueryRow(q, args...).Scan(&u.UID, &u.Role, &u.TPS, &u.Usr, &u.Pwd, &u.From, &u.To); err == sql.ErrNoRows {
		return u, new(InvalidCredentials)
	} else if err != nil {
		return u, err
	}

	pepper := config.Config().GetString("pepper")
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(password+pepper)); err != nil {
		return u, new(InvalidCredentials)
	}

	now := time.Now()
	if now.Before(u.From) || now.After(u.To) {
		return u, new(ExpiredCredentials)
	}

	return u, nil

}
