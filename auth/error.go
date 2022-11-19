package auth

import (
	"github.com/zicare/rgm/msg"
)

// AclTagsError exported
type AclTagsError struct {
	msg.Message
}

// UserTagsError exported
type UserTagsError struct {
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
