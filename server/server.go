package server

import (
	"log"
	"mahjong/handler"
	"net"

	"github.com/k-jun/northpole"
)

var (
	word1 = "wellcome to the cli-mahjong\n"
	word2 = "please enter your random user-id in uuid\nEnter> "
)

type Server interface {
	Run()
}

type serverImpl struct {
	listener net.Listener
	matches  *northpole.Match
}

func New(listener net.Listener) Server {
	m := northpole.New()

	return &serverImpl{
		listener: listener,
		matches:  &m,
	}
}

func (s *serverImpl) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		h := handler.New(conn, s.matches)
		go h.Run()
	}
}
