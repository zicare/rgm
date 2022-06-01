package db

//Table exported
type Table interface {
	New() Table
	Table() string
	View() string
	Value() interface{}
}
