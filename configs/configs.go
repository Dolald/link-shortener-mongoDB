package configs

import "time"

var (
	ContextWaiting time.Duration = 10
	Symbols                      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	Hours    time.Duration = 24
	IdLength               = 6
)
