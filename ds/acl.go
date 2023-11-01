package ds

import (
	"time"
)

// In-memory access control list.
// acl maps each grant to a time range.
// Helps speed up Authorization middleware.
var acl Acl

// Acl exported
type Acl map[Grant]TimeRange

//TimeRange exported
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// Grant exported
type Grant struct {
	Role   string `json:"role"`
	Route  string `json:"route"`
	Method string `json:"method"`
}

// Validates if g Grant exists and is valid at the time.
func (g Grant) Valid() bool {

	now := time.Now()
	if r, ok := acl[g]; !ok {
		return false
	} else if now.Before(r.From) || now.After(r.To) {
		return false
	}
	return true
}

// Defines an interface for ACL data access.
type IAclDataSource interface {

	// Returns all grants mapped to its corresponding validity time range.
	Fetch() (Acl, error)
}

// Meant to be executed on startup, Init loads the acl map in memory.
// acl maps each grant to a time range.
// Helps speed up Authorization middleware.
func Init(fn AclDSFactory, d IDataSource) (err error) {

	if dsrc, err := fn(d); err != nil {
		return err
	} else if acl, err = dsrc.Fetch(); err != nil {
		return err
	}

	return nil
}
