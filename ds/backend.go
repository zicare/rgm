package ds

import (
	"fmt"

	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

// Backend abstracts DB-scoped concerns for a single app-wide database.
type Backend interface {
	// Meta returns a registry with one entry per table/view, etc.
	// The map key is the real table/view name.
	Meta() (map[string]TableMeta, error)
}

type TableMeta struct {
	Columns []string
	PKs     []string
}

// global backend selection (one per app)
//var currentBackend Backend

// var backendSet bool
var backendType string
var backendMeta map[string]TableMeta

// must be called with magic init from any package
// that implements the Backend interface, like mysql, postgres, etc.
// this makes sure that only one backend type (mysql, postgres, etc.) is used
// and populates a registry with correct meta from the underlying database/backend.
func SetBackend(b Backend) (err error) {

	newType := fmt.Sprintf("%T", b)
	if backendType == "" {
		backendType = newType
	} else if newType != backendType {
		panic(fmt.Errorf("backend already set to %s and attempted to also set %s, only one backend type allowed", backendType, newType))
	}

	if backendMeta != nil {
		return nil // idempotent same-backend call
	}

	backendMeta, err = b.Meta()

	return err
}

func GetTableMeta(table string) (TableMeta, bool) {

	tm, ok := backendMeta[table]
	return tm, ok
}

//func backend() Backend { return currentBackend }

// ValidateTags ensures declared columns/PKs exist in the DB schema.
// Does not require that all DB PKs be tagged in the model.
func ValidateTags(m IDataSource, keys []string, fields []string) *TagError {

	tm, ok := GetTableMeta(m.Name())
	if !ok {
		return &TagError{
			Message: msg.Message{
				Msg: fmt.Sprintf("ValidateTags(%s): Table not found in registry", m.Name()),
			},
		}
	}

	// Declared columns must exist
	for _, col := range fields {
		if !lib.Contains(tm.Columns, col) {
			return &TagError{
				Message: msg.Message{
					Msg: fmt.Sprintf("ValidateTags(%s): %s is not a column", m.Name(), col),
				},
			}
		}
	}

	// If the table is a actually view skip PK validation
	if tm.PKs == nil {
		return nil
	}

	// If a real table, keys must be actual PKs
	for _, pk := range keys {
		if !lib.Contains(tm.PKs, pk) {
			return &TagError{
				Message: msg.Message{
					Msg: fmt.Sprintf("ValidateTags(%s): %s is not a primary key", m.Name(), pk),
				},
			}
		}
	}

	return nil
}
