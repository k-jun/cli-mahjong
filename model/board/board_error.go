package board

import "errors"

var (
	BoardMaxNOUErr             = errors.New("reach to the max number of users in the table")
	BoardIndexOutOfRangeErr    = errors.New("the index is out of range")
	BoardPlayerNotFoundErr     = errors.New("the player not found in the board")
	BoardActionAlreadyTokenErr = errors.New("requresed action is timeover")
)
