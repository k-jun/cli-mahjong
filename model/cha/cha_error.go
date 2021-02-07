package cha

import "errors"

var (
	ChaAlreadyHaveTumohaiErr = errors.New("already have tumohai")
	ChaAlreadyHaveYamaErr    = errors.New("already have yama")
	ChaAlreadyDidHaihaiErr   = errors.New("already did haihai")
	ChaHaiNotFoundErr        = errors.New("hai not found")
)
