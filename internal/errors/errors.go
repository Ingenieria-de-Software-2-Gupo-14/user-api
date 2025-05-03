package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrEmailInUsed = errors.New("user already exists")
