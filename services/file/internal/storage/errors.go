package storage

import "errors"

var ErrStorageDisabled = errors.New("object storage is disabled")
var ErrStorageInvalidInput = errors.New("storage input is invalid")
var ErrStorageObjectTooLarge = errors.New("storage object is too large")
