package domain

import "errors"

var (
	ErrNotFound      = errors.New("model was not found")
	ErrAlreadyExists = errors.New("model already exists")
)
