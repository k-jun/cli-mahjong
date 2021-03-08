package board

import (
	"mahjong/model/hai"
	"mahjong/model/player"
	"mahjong/model/yama"
	"sync"
)

var (
	MaxNumberOfUsers = 4
)

type Board interface {
	// getter
	Players() []*boardPlayer
	MaxNumberOfUser() int

	// setter
	SetWinIndex(int) error

	// game
	JoinPlayer(player.Player) (chan Board, error)
	LeavePlayer(player.Player) error
	Broadcast()

	// turn
	CurrentTurn() int
	NextTurn() int
	MyTurn(player.Player) (int, error)
	TurnEnd() error

	// last ho
	LastKawa() (*hai.Hai, error)

	// action counter
	ActionCounter() int
	CancelAction(c player.Player) error
	TakeAction(player.Player, func(*hai.Hai) error) error
}

func New(maxNOU int, y yama.Yama) Board {
	return &boardImpl{
		players:         []*boardPlayer{},
		yama:            y,
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		actionPlayers:   []*boardPlayer{},
		winIndex:        -1,
	}
}

type boardImpl struct {
	sync.Mutex
	players         []*boardPlayer
	yama            yama.Yama
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
	actionPlayers   []*boardPlayer

	// win
	winIndex int
}

type boardPlayer struct {
	channel chan Board
	player.Player
}

func (b *boardImpl) Players() []*boardPlayer {
	return b.players
}

func (b *boardImpl) MaxNumberOfUser() int {
	return b.maxNumberOfUser
}

func (t *boardImpl) SetWinIndex(idx int) error {
	if idx >= len(t.players) || idx < 0 {
		return BoardIndexOutOfRangeErr
	}
	t.winIndex = idx
	return nil
}

func (t *boardImpl) JoinPlayer(c player.Player) (chan Board, error) {
	t.Lock()
	defer t.Unlock()
	if len(t.players) >= t.maxNumberOfUser {
		return nil, BoardMaxNOUErr
	}

	if err := c.SetYama(t.yama); err != nil {
		return nil, err
	}
	channel := make(chan Board, t.maxNumberOfUser*3)
	t.players = append(t.players, &boardPlayer{Player: c, channel: channel})

	if len(t.players) >= t.maxNumberOfUser {
		t.gameStart()
		go t.Broadcast()
	}

	return channel, nil
}

func (t *boardImpl) LeavePlayer(c player.Player) error {
	t.Lock()
	defer t.Unlock()
	// terminate the game
	if t.isPlaying {
		t.isPlaying = false
		for _, tu := range t.players {
			close(tu.channel)
		}
		t.players = []*boardPlayer{}
	}
	return nil
}

func (t *boardImpl) Broadcast() {
	for _, tu := range t.players {
		tu.channel <- t
	}
}

func (t *boardImpl) gameStart() error {
	// tehai assign
	for _, tc := range t.players {
		if err := tc.Haipai(); err != nil {
			return err
		}
	}

	// tsumo
	return t.players[t.CurrentTurn()].Tsumo()
}

func (t *boardImpl) CurrentTurn() int {
	return t.turnIndex
}

func (t *boardImpl) MyTurn(c player.Player) (int, error) {
	for i, tc := range t.players {
		if tc.Player == c {
			return i, nil
		}
	}
	return -1, BoardPlayerNotFoundErr
}

func (t *boardImpl) NextTurn() int {
	return (t.turnIndex + 1) % t.maxNumberOfUser
}

func (t *boardImpl) TurnEnd() error {
	t.Lock()
	defer t.Unlock()
	err := t.setAction()
	if err != nil {
		return err
	}

	if len(t.actionPlayers) == 0 {
		if err := t.turnchange(t.NextTurn()); err != nil {
			return err
		}
		if err := t.players[t.CurrentTurn()].Tsumo(); err != nil {
			return err
		}
	}
	go t.Broadcast()
	return nil
}

func (t *boardImpl) turnchange(idx int) error {
	if idx < 0 || idx >= len(t.players) {
		return BoardIndexOutOfRangeErr
	}
	t.turnIndex = idx
	return nil
}

func (t *boardImpl) setAction() error {
	players := []*boardPlayer{}

	inHai, err := t.players[t.CurrentTurn()].Kawa().Last()
	if err != nil {
		return err
	}
	for i, tc := range t.players {
		if tc == t.players[t.CurrentTurn()] {
			continue
		}

		type Arg struct {
			ok bool
			e  error
		}
		args := []Arg{}
		flag := false
		if i == t.NextTurn() {
			ok, err := tc.Tehai().CanChii(inHai)
			args = append(args, Arg{ok, err})
		}
		ok, err := tc.Tehai().CanPon(inHai)
		args = append(args, Arg{ok, err})
		ok, err = tc.Tehai().CanMinKan(inHai)
		args = append(args, Arg{ok, err})
		ok, err = tc.Tehai().CanRon(inHai)
		args = append(args, Arg{ok, err})

		for _, arg := range args {
			if arg.e != nil {
				return arg.e
			}
			if arg.ok {
				flag = true
			}
		}

		if flag {
			players = append(players, tc)
		}
	}
	t.actionPlayers = players
	return nil
}

func (t *boardImpl) LastKawa() (*hai.Hai, error) {
	return t.players[t.CurrentTurn()].Kawa().Last()
}

func (t *boardImpl) ActionCounter() int {
	return len(t.actionPlayers)
}

func (t *boardImpl) CancelAction(c player.Player) error {
	t.Lock()
	defer t.Unlock()
	if len(t.actionPlayers) == 0 {
		return nil
	}

	found := false
	for i, tc := range t.actionPlayers {
		if tc.Player == c {
			found = true
			t.actionPlayers = append(t.actionPlayers[:i], t.actionPlayers[i+1:]...)
		}
	}
	if !found {
		return BoardPlayerNotFoundErr
	}

	if len(t.actionPlayers) == 0 {
		if err := t.turnchange(t.NextTurn()); err != nil {
			return err
		}
		if err := t.players[t.CurrentTurn()].Tsumo(); err != nil {
			return err
		}
		go t.Broadcast()
	}
	return nil
}

func (t *boardImpl) TakeAction(c player.Player, action func(*hai.Hai) error) error {
	t.Lock()
	defer t.Unlock()
	if len(t.actionPlayers) == 0 {
		return BoardActionAlreadyTokenErr
	}

	found := false
	for _, tc := range t.actionPlayers {
		if tc.Player == c {
			found = true
		}
	}
	if !found {
		return BoardPlayerNotFoundErr
	}

	h, err := t.players[t.CurrentTurn()].Kawa().Last()
	if err != nil {
		return err
	}
	if err := action(h); err != nil {
		return err
	}
	_, err = t.players[t.CurrentTurn()].Kawa().RemoveLast()
	if err != nil {
		return err
	}
	t.actionPlayers = []*boardPlayer{}

	turnIdx, err := t.MyTurn(c)
	if err != nil {
		return err
	}
	if err := t.turnchange(turnIdx); err != nil {
		return err
	}
	go t.Broadcast()
	return nil
}
