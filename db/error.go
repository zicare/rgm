package db

import (
	"github.com/zicare/rgm/msg"
)

//NotFoundError exported
type NotFoundError struct {
	msg.Message
}

//NotAllowedError exported
type NotAllowedError struct {
	msg.Message
}

//ConflictError exported
type ConflictError struct {
	msg.Message
}

//ParamError exported
type ParamError struct {
	msg.Message
}
