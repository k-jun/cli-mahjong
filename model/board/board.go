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

type ActionType string

var (
	Normal ActionType = "noaction"
	Tsumo  ActionType = "tsumo"
	Riichi ActionType = "riichi"
	Chii   ActionType = "chii"
	Pon    ActionType = "pon"
	Kan    ActionType = "kan"
	Ron    ActionType = "ron"
	Cancel ActionType = "no"
)

type Board interface {
	// getter
	Players() []*boardPlayer
	ActionPlayers() []*boardActionPlayer
	MaxNumberOfUser() int
	Winner() player.Player

	// setter
	SetWinner(player.Player) error

	// game
	JoinPlayer(player.Player) (chan Board, error)
	LeavePlayer(player.Player) error
	Broadcast()

	// turn
	CurrentTurn() int
	NextTurn() int
	MyTurn(player.Player) (int, error)
	TurnEnd() error

	// last hai
	LastKawa() (*hai.Hai, error)

	// actions
	MyAction(p player.Player) ([]ActionType, error)
	CancelAction(c player.Player) error
	TakeAction(player.Player, func(*hai.Hai) error) error
}

func New(maxNOU int, y yama.Yama) Board {
	return &boardImpl{
		players:         []*boardPlayer{},
		actionPlayers:   []*boardActionPlayer{},
		yama:            y,
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		winner:          nil,
	}
}

type boardImpl struct {
	sync.Mutex
	players         []*boardPlayer
	actionPlayers   []*boardActionPlayer
	yama            yama.Yama
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool

	// win
	winner player.Player
}

type boardPlayer struct {
	channel chan Board
	player.Player
}

type boardActionPlayer struct {
	actions []ActionType
	player.Player
}

func (b *boardImpl) Players() []*boardPlayer {
	return b.players
}

func (t *boardImpl) ActionPlayers() []*boardActionPlayer {
	return t.actionPlayers
}

func (b *boardImpl) MaxNumberOfUser() int {
	return b.maxNumberOfUser
}

func (b *boardImpl) Winner() player.Player {
	return b.winner
}

func (t *boardImpl) SetWinner(p player.Player) error {
	if p == nil {
		return BoardPlayerNilError
	}
	t.winner = p
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
	err := t.setActionPlayer()
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

func (t *boardImpl) setActionPlayer() error {
	actionPlayers := []*boardActionPlayer{}

	inHai, err := t.players[t.CurrentTurn()].Kawa().Last()
	if err != nil {
		return err
	}
	for i, tc := range t.players {
		if tc == t.players[t.CurrentTurn()] {
			continue
		}

		type Arg struct {
			ok     bool
			e      error
			action ActionType
		}
		args := []Arg{}
		if i == t.NextTurn() {
			ok, err := tc.CanChii(inHai)
			args = append(args, Arg{ok, err, Chii})
		}
		ok, err := tc.CanPon(inHai)
		args = append(args, Arg{ok, err, Pon})
		ok, err = tc.CanMinKan(inHai)
		args = append(args, Arg{ok, err, Kan})
		ok, err = tc.CanRon(inHai)
		args = append(args, Arg{ok, err, Ron})

		actions := []ActionType{}
		for _, arg := range args {
			if arg.e != nil {
				return arg.e
			}
			if arg.ok {
				actions = append(actions, arg.action)
			}
		}
		if len(actions) != 0 {
			actionPlayer := boardActionPlayer{
				Player:  tc.Player,
				actions: actions,
			}
			actionPlayers = append(actionPlayers, &actionPlayer)

		}
	}
	t.actionPlayers = actionPlayers
	return nil
}

func (t *boardImpl) LastKawa() (*hai.Hai, error) {
	return t.players[t.CurrentTurn()].Kawa().Last()
}

func (t *boardImpl) MyAction(p player.Player) ([]ActionType, error) {
	for _, ap := range t.actionPlayers {
		if ap.Player == p {
			return ap.actions, nil
		}
	}
	return []ActionType{}, nil
}

func (t *boardImpl) CancelAction(p player.Player) error {
	t.Lock()
	defer t.Unlock()
	if len(t.actionPlayers) == 0 {
		return nil
	}

	found := false
	for i, tc := range t.actionPlayers {
		if tc.Player == p {
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
	t.actionPlayers = []*boardActionPlayer{}

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
