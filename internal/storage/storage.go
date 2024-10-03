package storage

import "errors"

var (
	ErrGuitarNotFound = errors.New("guitar not found")
	ErrGuitarExists   = errors.New("guitar exists")
)
