package mysql

import "github.com/zicare/rgm/msg"

// NotITableError exported
type NotITableError struct {
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
