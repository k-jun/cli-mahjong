package storage

import "errors"

var (
	TakuStorageAlreadyExistErr = errors.New("a taku having the id already exist")
	TakuStorageNotExistErr     = errors.New("a taku having the id not exist")
)
