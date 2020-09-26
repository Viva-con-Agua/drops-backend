package utils

import "errors"

var (
	ErrorUserNotFound     = errors.New("user not found")
	ErrorUserNotConfirmed = errors.New("user not confirmed")
)
