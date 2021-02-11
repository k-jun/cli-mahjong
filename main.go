package main

import (
	"log"
	"mahjong/server"
	"net"
)

// func main() {
// 	tk := taku.New(taku.MaxNumberOfUsers)
// 	h := &ho.HoMock{HaisMock: append(hai.Manzu)}
// 	huroMock := &huro.HuroMock{ChiisMock: [][3]*hai.Hai{{hai.Chun, hai.Chun, hai.Chun}}}
// 	c := cha.New(utils.NewUUID(), h, tehai.New(), nil, huroMock)
// 	_, err := tk.JoinCha(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < taku.MaxNumberOfUsers-1; i++ {
// 		t := tehai.New()
// 		hu := huro.New()
// 		ho := &ho.HoMock{HaisMock: append(hai.Manzu)}
// 		cha := cha.New(utils.NewUUID(), ho, t, nil, hu)
// 		tk.JoinCha(cha)
// 	}
// 	c.Tsumo()
// 	fmt.Println(tk.Draw(c))
//
// 	// fmt.Println(str)
// }

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(ln)
	s.Run()
}
