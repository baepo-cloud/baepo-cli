package baepoerrors

import "errors"

var (
	AuthError        = errors.New("authentication error")
	ConfigError      = errors.New("configuration error")
	MachineError     = errors.New("machine error")
	InvalidArgsError = errors.New("invalid arguments")
)
