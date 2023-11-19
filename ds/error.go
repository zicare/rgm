package ds

import (
	"encoding/json"

	"github.com/zicare/rgm/msg"
)

// ValidationErrors exported
type ValidationErrors msg.MessageList

//Error exported
func (ve *ValidationErrors) Error() string {

	jve, _ := json.Marshal(ve)
	return string(jve)
}

// ValidationError exported
type ValidationError struct {
	msg.Message
}

// TagError exported
type TagError struct {
	msg.Message
}

// DuplicatedEntry exported
type DuplicatedEntry struct {
	msg.Message
}

// ForeignKeyConstraint exported
type ForeignKeyConstraint struct {
	msg.Message
}

// NotFoundError exported
type NotFoundError struct {
	msg.Message
}

// NotAllowedError exported
type NotAllowedError struct {
	msg.Message
}

// InvalidCredentials exported
type InvalidCredentials struct {
	msg.Message
}

// ExpiredCredentials exported
type ExpiredCredentials struct {
	msg.Message
}

// InvalidPIN exported
type InvalidPinError struct {
	msg.Message
}

// ExpiredPIN exported
type ExpiredPinError struct {
	msg.Message
}

// InsertError exported
type InsertError struct {
	msg.Message
}

// UpdateError exported
type UpdateError struct {
	msg.Message
}
