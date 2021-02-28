package cha

import "errors"

var (
	ChaAlreadyHaveTsumohaiErr = errors.New("already have tsumohai")
	ChaAlreadyHaveYamaErr     = errors.New("already have yama")
	ChaAlreadyDidHaipaiErr    = errors.New("already did haipai")
	ChaHaiNotFoundErr         = errors.New("hai not found")
	ChaAlreadyRiichiErr       = errors.New("already did riichi")
)
