package handler

import (
	"fmt"
	"log"
	"mahjong/model/kawa"
	"mahjong/model/naki"
	"mahjong/model/player"
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
	fmt.Println("roomId: ", roomId)
	if err != nil {
		log.Println(err)
		return
	}

	t := tehai.New()
	n := naki.New()
	k := kawa.New()
	cha := player.New(h.id, k, t, n)
	roomChan, err := h.gameUsecase.JoinBoard(roomId, cha)
	if err != nil {
		log.Println(err)
		return
	}
	go h.gameUsecase.InputController(roomId, cha)
	err = h.gameUsecase.OutputController(roomId, cha, roomChan)
	if err != nil {
		log.Println(err)
	}
}
