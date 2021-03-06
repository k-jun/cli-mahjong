package usecase

import "errors"

var (
	MatchUsecaseRoomChannelClosedErr = errors.New("the room channel closed")
	GameUsecaseBoardChannelClosedErr = errors.New("the board channel closed")
	GameUsecaseInvalidActionErr      = errors.New("invalid action")
)
