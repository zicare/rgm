package auth

import (
	"encoding/json"
	"time"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
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
		return u, new(UserTagsError)
	}

	w := make(db.UParams)
	w[f[3]] = username

	if _, data, err := db.Find(db.QueryOptionsFactory(ds.t, "", nil, w)); err != nil {
		return u, err
	} else {
		data, _ := json.Marshal(data)
		json.Unmarshal(data, &u)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(password+pepper)); err != nil {
		return u, new(InvalidCredentials)
	}

	if now.Before(u.From) || now.After(u.To) {
		return u, new(ExpiredCredentials)
	}

	return u, nil

}
