package db

// TableMeta exported
type TableMeta struct {
	Fields   []string
	Primary  []string
	Serial   []string
	View     []string
	Writable []string
}

// Table exported
type Table interface {

	// Must return the table name
	Name() string

	// Must attach foreign table data if available
	Dig()

	// Must set conditions to filter out content
	// not intended for the user making the request
	Scope(uid string, t ...Table) map[string]string
}

type BaseTable struct{}

// Scope exported
func (BaseTable) Scope(uid string, t ...Table) map[string]string {

	return make(map[string]string)
}

// Dig exported
func (BaseTable) Dig() {}
