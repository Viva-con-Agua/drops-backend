package utils

import "errors"

var (
	ErrorNotFound = errors.New("NotFound")
	ErrorConflict = errors.New("Conflict")
	ErrorPassword = errors.New("Password")
)
