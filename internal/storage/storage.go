package storage

import "errors"

var (
	ErrURLExists     = errors.New("url already exists")
	ErrAliasExists   = errors.New("alias already exists")
	ErrURLNotFound   = errors.New("url not found")
	ErrAliasNotFound = errors.New("alias not found")
)
