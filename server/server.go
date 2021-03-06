package server

import (
	"log"
	"mahjong/model/board"
	"mahjong/model/yama"
	"mahjong/server/handler"
	"mahjong/server/usecase"
	"mahjong/storage"
	"mahjong/utils"
	"net"

	"github.com/k-jun/northpole"
)

type Server interface {
	Run()
}

type serverImpl struct {
	listener     net.Listener
	matches      northpole.Match
	boardStorage storage.BoardStorage
}

func New(listener net.Listener) Server {
	m := northpole.New()
	ts := storage.NewBoardStorage()

	return &serverImpl{
		listener:     listener,
		matches:      m,
		boardStorage: ts,
	}
}

func (s *serverImpl) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		write := func(mess string) error {
			_, err := conn.Write([]byte(mess))
			return err
		}
		read := func(buffer []byte) error {
			_, err := conn.Read(buffer)
			return err
		}

		callback := func(id string) error {
			yama := yama.New()
			taku := board.New(board.MaxNumberOfUsers, yama)
			s.boardStorage.Add(id, taku)
			return nil
		}
		close := func() error {
			return conn.Close()
		}

		matchUsecase := usecase.NewMatchUsecase(s.matches, write, read, callback)
		gameUsecase := usecase.NewGameUsecase(s.boardStorage, write, read)
		id := utils.NewUUID()
		h := handler.New(id, matchUsecase, gameUsecase, close)
		go h.Run()
	}
}
