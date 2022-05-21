package msg

import (
	"fmt"
)

//Message exported
type Message struct {
	Key   string
	Msg   string
	Args  []interface{}
	Field string
}

//SetArgs exported
func (m Message) SetArgs(args ...interface{}) Message {
	m.Args = args
	return m
}

//SetField exported
func (m Message) SetField(field string) Message {
	m.Field = field
	return m
}

//M2E exported
func (m Message) M2E() *Message {
	return &m
}

//Error exported
func (m *Message) Error() string {
	return m.String()
}

//String exported
func (m Message) String() string {

	if m.Args != nil && len(m.Args) > 0 {
		return fmt.Sprintf(m.Msg, m.Args)
	}
	return m.Msg
}
