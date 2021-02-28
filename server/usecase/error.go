package usecase

import "errors"

var (
	MatchUsecaseRoomChannelClosedErr = errors.New("the room channel closed")
	GameUsecaseTakuChannelClosedErr  = errors.New("the taku channel closed")
	GameUsecaseInvalidActionErr      = errors.New("invalid action")
)
