package db

//View exported
type View struct{}

//Table exported
func (View) Table() string {
	return ""
}
