package taku

import (
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/yama"
	"sync"
)

var (
	MaxNumberOfUsers = 4
)

type Taku interface {
	// game
	JoinCha(cha.Cha) (chan Taku, error)
	LeaveCha(cha.Cha) error
	Broadcast()

	// turn
	CurrentTurn() int
	MyTurn(cha.Cha) int
	TurnEnd() error

	// last ho
	LastHo() (*hai.Hai, error)

	// action counter
	CancelAction() error
	TakeAction(func(*hai.Hai) error) error
}

func New(maxNOU int) Taku {
	return &takuImpl{
		chas:            []*takuCha{},
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		actionCounter:   0,
	}
}

type takuImpl struct {
	sync.Mutex
	chas            []*takuCha
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
	actionCounter   int
}

type takuCha struct {
	channel chan Taku
	cha     cha.Cha
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

func (t *takuImpl) Broadcast() {
	for _, tu := range t.chas {
		tu.channel <- t
	}
}

func (t *takuImpl) gameStart() error {
	// create yama
	y := yama.New()
	if err := y.Kan(); err != nil {
		return err
	}

	// tehai assign
	for _, tc := range t.chas {
		if err := tc.cha.SetYama(y); err != nil {
			return err
		}
		if err := tc.cha.Haipai(); err != nil {
			return err
		}
	}

	return nil
}

func (t *takuImpl) CurrentTurn() int {
	return t.turnIndex
}

func (t *takuImpl) MyTurn(c cha.Cha) int {
	for i, tc := range t.chas {
		if tc.cha == c {
			return i
		}
	}
	return -1
}

func (t *takuImpl) nextTurn() int {
	return (t.turnIndex + 1) % t.maxNumberOfUser
}

func (t *takuImpl) TurnEnd() error {
	t.Lock()
	defer t.Unlock()
	err := t.setActionCounter()
	if err != nil {
		return err
	}

	if t.actionCounter == 0 {
		if err := t.turnChange(t.nextTurn()); err != nil {
			return err
		}
	}
	go t.Broadcast()
	return nil
}

func (t *takuImpl) turnChange(idx int) error {
	if idx < 0 || idx >= len(t.chas) {
		return TakuIndexOutOfRangeErr
	}
	t.turnIndex = idx
	return nil
}

func (t *takuImpl) setActionCounter() error {
	counter := 0

	inHai, err := t.chas[t.CurrentTurn()].cha.Ho().Last()
	if err != nil {
		return err
	}
	for _, tc := range t.chas {
		if tc == t.chas[t.CurrentTurn()] {
			continue
		}
		actions, err := tc.cha.FindHuroActions(inHai)
		if err != nil {
			return err
		}
		if len(actions) != 0 {
			counter += 1
		}
	}
	t.actionCounter = counter
	return nil
}

func (t *takuImpl) LastHo() (*hai.Hai, error) {
	return t.chas[t.CurrentTurn()].cha.Ho().Last()
}

func (t *takuImpl) CancelAction() error {
	t.Lock()
	defer t.Unlock()
	if t.actionCounter == 0 {
		return nil
	}
	t.actionCounter -= 1
	if t.actionCounter == 0 {
		if err := t.turnChange(t.nextTurn()); err != nil {
			return err
		}
		go t.Broadcast()
	}
	return nil
}

func (t *takuImpl) TakeAction(action func(*hai.Hai) error) error {
	t.Lock()
	defer t.Unlock()
	if t.actionCounter == 0 {
		return TakuActionAlreadyTokenErr
	}
	t.actionCounter = 0
	h, err := t.chas[t.CurrentTurn()].cha.Ho().RemoveLast()
	if err != nil {
		return err
	}
	if err := action(h); err != nil {
		return err
	}
	go t.Broadcast()
	return nil
}
