package player

import "errors"

var (
	PlayerAlreadyHaveTsumohaiErr = errors.New("already have tsumohai")
	PlayerAlreadyHaveYamaErr     = errors.New("already have yama")
	PlayerAlreadyDidHaipaiErr    = errors.New("already did haipai")
	PlayerHaiNotFoundErr         = errors.New("hai not found")
	PlayerAlreadyRiichiErr       = errors.New("already did riichi")
	PlayerActionInvalidErr       = errors.New("invalid Player action")
)
