package jwt

import (
	"github.com/zicare/rgm/msg"
)

// InvalidToken exported
type InvalidToken struct {
	msg.Message
}

// InvalidTokenPayload exported
type InvalidTokenPayload struct {
	msg.Message
}

// TamperedToken exported
type TamperedToken struct {
	msg.Message
}

// ExpiredToken exported
type ExpiredToken struct {
	msg.Message
}
