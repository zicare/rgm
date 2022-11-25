package db

import (
	"github.com/zicare/rgm/msg"
)

// OpenConnError exported
type OpenConnError struct {
	msg.Message
}

// PingTestError exported
type PingTestError struct {
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

// ConflictError exported
type ConflictError struct {
	msg.Message
}

// ParamError exported
type ParamError struct {
	msg.Message
}

// TableTagError exported
type TableTagError struct {
	msg.Message
}
