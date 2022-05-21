package msg

import (
	"encoding/json"
)

//MessageList exported
type MessageList []Message

//Error exported
func (ml MessageList) Error() string {

	jml, _ := json.Marshal(ml)
	return string(jml)
}
