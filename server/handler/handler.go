package handler

import (
	"mahjong/server/usecase"

	"github.com/google/uuid"
	"github.com/k-jun/northpole/user"
)

type Handler interface {
	Run()
}

var MaxNumberOfUsers = 4

type handlerImpl struct {
	id           uuid.UUID
	matchUsecase usecase.MatchUsecase
	close        func() error
}

func New(id uuid.UUID, matchUsecase usecase.MatchUsecase, close func() error) Handler {
	return &handlerImpl{id, matchUsecase, close}
}

func (h *handlerImpl) Run() {
	defer h.close()
	user := user.New(h.id)

	_, err := h.matchUsecase.JoinRandomRoom(user)
	if err != nil {
		return
	}
}
