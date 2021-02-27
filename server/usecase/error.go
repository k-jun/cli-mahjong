package usecase

import "errors"

var (
	MatchUsecaseRoomChannelClosedErr = errors.New("the room channel closed")
	GameUsecaseTakuChannelClosedErr  = errors.New("the taku channel closed")
	GameUsecaseHuroNotFoundErr       = errors.New("the huro action not found")
)
