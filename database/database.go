package database

import (
	"github.com/pkg/errors"
)

// ErrNotFound is returned if there is no key in DB
var ErrNotFound = errors.New("Not found")
