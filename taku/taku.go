package taku

import (
	"mahjong/cha"
	"sync"
)

type Taku interface {
	JoinCha(cha.Cha) (chan Taku, error)
	LeaveCha(cha.Cha) error
	Broadcast()
	NextTurn(int) error
}

func New(maxNOU int) Taku {
	return &takuImpl{
		chas:            []*takuCha{},
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       false,
	}
}

type takuImpl struct {
	sync.Mutex
	chas            []*takuCha
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
}

type takuCha struct {
	channel chan Taku
	cha     cha.Cha
}

func (t *takuImpl) JoinCha(c cha.Cha) (chan Taku, error) {
	if len(t.chas) >= t.maxNumberOfUser {
		return nil, TakuMaxNOUErr
	}
	channel := make(chan Taku, t.maxNumberOfUser*3)
	t.chas = append(t.chas, &takuCha{cha: c, channel: channel})
	if len(t.chas) >= t.maxNumberOfUser {
		t.isPlaying = true

	}
	go t.Broadcast()

	return channel, nil
}

func (t *takuImpl) LeaveCha(_ cha.Cha) error {
	// terminate the game
	for _, tu := range t.chas {
		close(tu.channel)
	}
	return nil
}

func (t *takuImpl) NextTurn(idx int) error {
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
