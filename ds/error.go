package ds

import "github.com/zicare/rgm/msg"

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
