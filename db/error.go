package db

import (
	"github.com/zicare/rgm/msg"
)

//NotFoundError exported
type NotFoundError struct {
	Message msg.Message
}

//Error exported
func (e *NotFoundError) Error() string {

	return e.Message.Error()
}

//NotAllowedError exported
type NotAllowedError struct {
	Message msg.Message
}

//Error exported
func (e *NotAllowedError) Error() string {

	return e.Message.Error()
}

//ConflictError exported
type ConflictError struct {
	Message msg.Message
}

//Error exported
func (e *ConflictError) Error() string {

	return e.Message.Error()
}

//ParamError exported
type ParamError struct {
	Message msg.Message
}

//Error exported
func (e *ParamError) Error() string {

	return e.Message.Error()
}
