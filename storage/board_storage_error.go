package storage

import "errors"

var (
	BoardStorageAlreadyExistErr = errors.New("a board having the id already exist")
	BoardStorageNotExistErr     = errors.New("a board having the id not exist")
)
