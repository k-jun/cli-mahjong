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

// func handleConnection(conn net.Conn) {
// 	bytes := []byte("wellcome to the cli-mahjong\n")
// 	conn.Write(bytes)
// 	buffer := make([]byte, 1024)
// 	_, err := conn.Read(buffer)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(buffer))
// }
