package msgr

import "errors"

var (
	ErrInvalidMessage = errors.New("invalid message")
	ErrNoProviders    = errors.New("no providers found")
	ErrInvalidFormat  = errors.New(`invalid format, needs to be "html" or "text"`)
)
