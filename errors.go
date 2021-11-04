package it

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidPath = errors.New("invalid path")
	ErrInvalidName = errors.New("invalid name")
	ErrBadOrder    = errors.New("bad order")
)
