package auth

import (
	"encoding/json"
	"time"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
	"golang.org/x/crypto/bcrypt"
)

// UserDS implementations are used
// to authenticate users by username and password.
type UserDS interface {

	// Return User matching the username and password
	GetUser(username, password string) (User, error)
}

// Default implementation of UserDS interface.
type TUserDS struct {
	t db.Table
}

// TUserDSFactory exported
func TUserDSFactory(t db.Table) TUserDS {

	var ds TUserDS

	ds.t = t
	return ds
}

// GetUser exported
func (ds TUserDS) GetUser(username, password string) (User, error) {

	var (
		u      = User{Type: ds.t.Name()}
		now    = time.Now()
		pepper = config.Config().GetString("pepper")
	)

	// Verify user tags
	f, err := db.TaggedFields(ds.t, "json", []string{"uid", "role", "tps", "usr", "pwd", "from", "to"})
	if err != nil {
		// User tags are not properly set
		e := UserTagsError{Message: msg.Get("2").SetArgs("User")}
		return u, &e
	}

	w := make(db.Params)
	w[f[3]] = username

	if fo, err := db.FindOptionsFactory(ds.t, "", nil, w, false); err != nil {
		return u, err
	} else if _, data, err := db.Find(fo); err != nil {
		return u, err
	} else {
		data, _ := json.Marshal(data)
		json.Unmarshal(data, &u)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(password+pepper))
	switch err {
	case nil:
		//pwd okay
	case bcrypt.ErrMismatchedHashAndPassword:
		// Wrong password
		e := InvalidCredentials{Message: msg.Get("4")}
		return u, &e
	default:
		// Something went wrong.
		// Possibly a mal formed hashed password in database.
		e := InvalidCredentials{Message: msg.Get("4")}
		return u, &e
	}

	if now.Before(u.From) || now.After(u.To) {
		// Credentials expired or not yet valid
		e := ExpiredCredentials{Message: msg.Get("6")}
		return u, &e
	}

	return u, nil

}
