package main

import (
	"log"
	"mahjong/server"
	"net"
	"os"
)

// func main() {
// 	y := yama.New()
// 	tk := taku.New(taku.MaxNumberOfUsers, y)
// 	h := &ho.HoMock{HaisMock: append(hai.Manzu)}
// 	huroMock := &huro.HuroMock{ChiisMock: [][3]*hai.Hai{{hai.Chun, hai.Chun, hai.Chun}}}
// 	c := cha.New(utils.NewUUID(), h, tehai.New(), huroMock)
// 	_, err := tk.JoinCha(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < taku.MaxNumberOfUsers-1; i++ {
// 		t := tehai.New()
// 		hu := huro.New()
// 		ho := &ho.HoMock{HaisMock: append(hai.Manzu)}
// 		cha := cha.New(utils.NewUUID(), ho, t, hu)
// 		cha.Riichi(nil)
// 		tk.JoinCha(cha)
// 	}
// 	c.Tsumo()
// 	fmt.Println(tk.Draw(c))
//
// 	// fmt.Println(str)
// }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(ln)
	s.Run()
}
