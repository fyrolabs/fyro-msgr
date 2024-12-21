package msgr

import "errors"

var (
	ErrNoChannels    = errors.New("no channels found")
	ErrInvalidKind   = errors.New("invalid entry")
	ErrNoProviders   = errors.New("no providers found")
	ErrInvalidFormat = errors.New(`invalid format, needs to be "html" or "text"`)
	ErrNoLocaleFiles = errors.New("no locale files found")
)
