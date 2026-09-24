package streams

import "errors"

var (
	errInvalidUrl       = errors.New("invalid stream url")
	errInvalidStreamKey = errors.New("invalid stream key")
	errBanned           = errors.New("user is banned")
	errChannelMismatch  = errors.New("stream key does not belong to the channel")
)
