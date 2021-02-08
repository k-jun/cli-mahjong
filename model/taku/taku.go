package taku

import (
	"fmt"
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/yama"
	"sync"
)

var (
	MaxNumberOfUsers = 4
)

type Taku interface {
	JoinCha(cha.Cha) (chan Taku, error)
	LeaveCha(cha.Cha) error
	Broadcast()
	TurnChange(int) error
	CurrentTurn() int
	NextTurn() int
	IsYourTurn(cha.Cha) bool
	HasChaActions() (bool, error)
	LastDahai() (*hai.Hai, error)
	ChaActionCnt() int
}

func New(maxNOU int) Taku {
	return &takuImpl{
		chas:            []*takuCha{},
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		chaActionCnt:    0,
	}
}

type takuImpl struct {
	sync.Mutex
	chas            []*takuCha
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
	chaActionCnt    int
}

type takuCha struct {
	channel chan Taku
	cha     cha.Cha
}

func (t *takuImpl) IsYourTurn(c cha.Cha) bool {
	return t.chas[t.turnIndex].cha == c
}

func (t *takuImpl) ChaActionCnt() int {
	return t.chaActionCnt
}

func (t *takuImpl) JoinCha(c cha.Cha) (chan Taku, error) {
	t.Lock()
	defer t.Unlock()
	if len(t.chas) >= t.maxNumberOfUser {
		return nil, TakuMaxNOUErr
	}
	channel := make(chan Taku, t.maxNumberOfUser*3)
	t.chas = append(t.chas, &takuCha{cha: c, channel: channel})

	if len(t.chas) >= t.maxNumberOfUser {
		t.gameStart()
		go t.Broadcast()
	}

	return channel, nil
}

func (t *takuImpl) LeaveCha(c cha.Cha) error {
	t.Lock()
	defer t.Unlock()
	// terminate the game
	if t.isPlaying {
		t.isPlaying = false
		for _, tu := range t.chas {
			close(tu.channel)
		}
		t.chas = []*takuCha{}
	}
	return nil
}

func (t *takuImpl) CurrentTurn() int {
	return t.turnIndex
}

func (t *takuImpl) NextTurn() int {
	return (t.turnIndex + 1) % t.maxNumberOfUser
}

func (t *takuImpl) HasChaActions() (bool, error) {
	chaActionCnt := 0

	inHai, err := t.chas[t.CurrentTurn()].cha.Ho().Last()
	fmt.Println("inHai:", inHai)
	if err != nil {
		return false, err
	}
	for _, tc := range t.chas {
		if tc != t.chas[t.CurrentTurn()] {
			actions := tc.cha.CanHuro(inHai)
			fmt.Println("actions:", actions)
			if len(actions) != 0 {
				chaActionCnt += 1
			}
		}
	}
	t.chaActionCnt = chaActionCnt
	return chaActionCnt == 0, nil
}

func (t *takuImpl) LastDahai() (*hai.Hai, error) {
	return t.chas[t.CurrentTurn()].cha.Ho().Last()
}

func (t *takuImpl) TurnChange(idx int) error {
	t.Lock()
	defer t.Unlock()
	if idx < 0 || idx >= len(t.chas) {
		return TakuIndexOutOfRangeErr
	}
	t.turnIndex = idx
	go t.Broadcast()
	return nil
}

func (t *takuImpl) Broadcast() {
	for _, tu := range t.chas {
		tu.channel <- t
	}
}

func (t *takuImpl) gameStart() error {
	// create yama
	y := yama.New()
	if err := y.KanDora(); err != nil {
		return err
	}

	// tehai assign
	for _, tc := range t.chas {
		if err := tc.cha.SetYama(y); err != nil {
			return err
		}
		if err := tc.cha.Haihai(); err != nil {
			return err
		}
	}

	return nil
}
