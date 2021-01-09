package main

import (
	"log"
	"mahjong/server"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(ln)
	s.Run()
}
