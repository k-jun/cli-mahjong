package handler

import (
	"net"

	"github.com/k-jun/northpole"
)

type Handler interface {
	Run()
}

type handlerImpl struct {
	conn    net.Conn
	matches *northpole.Match
}

func New(conn net.Conn, matches *northpole.Match) Handler {
	return &handlerImpl{conn, matches}

}

func (h *handlerImpl) Run() {
	h.conn.Write([]byte(word1))

}
