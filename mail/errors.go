package mail

import "errors"

var (
	ErrInvalidKind     = errors.New("invalid entry")
	ErrInvalidProvider = errors.New("invalid send provider")
	ErrNoLocaleFiles   = errors.New("no locale files found")
)
