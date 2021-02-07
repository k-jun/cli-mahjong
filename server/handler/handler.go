package handler

import (
	"mahjong/model/cha"
	"mahjong/model/ho"
	"mahjong/model/huro"
	"mahjong/model/tehai"
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
	gameUsecase  usecase.GameUsecase
	close        func() error
}

func New(id uuid.UUID, matchUsecase usecase.MatchUsecase, gameUsecase usecase.GameUsecase, close func() error) Handler {
	return &handlerImpl{id, matchUsecase, gameUsecase, close}
}

func (h *handlerImpl) Run() {
	defer h.close()
	user := user.New(h.id.String())

	roomId, err := h.matchUsecase.JoinRandomRoom(user)
	if err != nil {
		return
	}

	t := tehai.New()
	hu := huro.New()
	ho := ho.New()
	cha := cha.New(h.id, ho, t, nil, hu)
	roomChan, err := h.gameUsecase.JoinTaku(roomId, cha)
	go h.gameUsecase.InputController(roomId, cha)
	h.gameUsecase.OutputController(roomId, cha, roomChan)
	if err != nil {
		return
	}
}
