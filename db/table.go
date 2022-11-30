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
	Dig(f ...string) []Dig

	// Must set conditions to filter out content
	// not intended for the user uid making the request
	// or not a child of optional parent t records
	Scope(uid string, t ...Table) map[string]string

	// Must set conditions to filter out content
	// not intended for the user uid making the request
	// or not a child of optional parent t records.
	// Second return parameter indicates whether or not
	// is okay to proceed with delete action.
	BeforeDelete(uid string, t ...Table) (map[string]string, bool)
}

type BaseTable struct{}

// Dig exported
func (BaseTable) Dig(f ...string) []Dig {

	return []Dig{}
}

// Scope exported
func (BaseTable) Scope(uid string, t ...Table) map[string]string {

	return make(map[string]string)
}

// Scope exported
func (BaseTable) BeforeDelete(uid string, t ...Table) (map[string]string, bool) {

	return make(map[string]string), true
}
