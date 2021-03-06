package board

import (
	"mahjong/model/hai"
	"mahjong/model/player"
	"mahjong/model/yama"
	"strings"
	"sync"
)

var (
	MaxNumberOfUsers = 4
)

type Board interface {
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

	// draw
	Draw(player.Player) string
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
	err := t.setActionCounter()
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

func (t *boardImpl) setActionCounter() error {
	players := []*boardPlayer{}

	inHai, err := t.players[t.CurrentTurn()].Kawa().Last()
	if err != nil {
		return err
	}
	for i, tc := range t.players {
		if tc == t.players[t.CurrentTurn()] {
			continue
		}

		actionCounter := 0
		if i == t.NextTurn() {
			pairs, err := tc.Tehai().FindChiiPairs(inHai)
			if err != nil {
				return err
			}
			actionCounter += len(pairs)
		}
		pairs, err := tc.Tehai().FindPonPairs(inHai)
		if err != nil {
			return err
		}
		actionCounter += len(pairs)
		kpairs, err := tc.Tehai().FindKanPairs(inHai)
		if err != nil {
			return err
		}
		actionCounter += len(kpairs)
		ok, err := tc.CanRon(inHai)
		if err != nil {
			return err
		}
		if ok {
			actionCounter += 1
		}

		if actionCounter != 0 {
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

	t.actionPlayers = []*boardPlayer{}
	h, err := t.players[t.CurrentTurn()].Kawa().RemoveLast()
	if err != nil {
		return err
	}

	if err := action(h); err != nil {
		return err
	}

	turnIdx, _ := t.MyTurn(c)
	if err := t.turnchange(turnIdx); err != nil {
		return err
	}
	go t.Broadcast()
	return nil
}

type boardHai struct {
	*hai.Hai
	isOpen   bool
	isDown   bool
	isRiichi bool
}

func (t *boardImpl) Draw(c player.Player) string {

	str := ""
	if t.winIndex != -1 {
		draftTehai := t.draftTehai(t.players[t.winIndex])
		str += drawTehai(draftTehai)
		str += "\n GAME SET!!"
		return str
	}
	// tehais
	draftTehais := t.draftTehaiAll(c)
	str += drawTehai(draftTehais["toimen"])

	// ho
	draftHo := t.draftHo(c)

	for i, h := range draftTehais["kamiplayer"] {
		if h == nil {
			continue
		}
		draftHo[i][0] = h
	}
	for i, h := range draftTehais["shimoplayer"] {
		if h == nil {
			continue
		}
		draftHo[len(draftHo)-1-i][len(draftHo[i])-1] = h
	}
	str += drawHo(draftHo)

	// tehai
	str += drawTehai(draftTehais["jiplayer"])

	return str
}

func (t *boardImpl) draftTehaiAll(c player.Player) map[string][20]*boardHai {
	tehaiMap := map[string][20]*boardHai{}
	idx, err := t.MyTurn(c)
	if err != nil {
		panic(err)
	}
	jiplayer := t.players[idx]
	shimoplayer := t.players[(idx+1)%t.maxNumberOfUser]
	toimen := t.players[(idx+2)%t.maxNumberOfUser]
	kamiplayer := t.players[(idx+3)%t.maxNumberOfUser]
	tehaiMap["jiplayer"] = t.draftTehai(jiplayer)
	tehaiMap["shimoplayer"] = hideTehai(t.draftTehai(shimoplayer))
	tehaiMap["toimen"] = hideTehai(reverse(t.draftTehai(toimen)))
	tehaiMap["kamiplayer"] = hideTehai(t.draftTehai(kamiplayer))
	return tehaiMap
}

func reverse(s [20]*boardHai) [20]*boardHai {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func hideTehai(v [20]*boardHai) [20]*boardHai {
	hais := [20]*boardHai{}
	for i, h := range v {
		if h != nil && !h.isDown {
			h.isOpen = false
		}
		hais[i] = h
	}

	return hais
}

func (t *boardImpl) draftTehai(c player.Player) [20]*boardHai {
	hais := [20]*boardHai{}
	for i, h := range c.Tehai().Hais() {
		hais[i] = &boardHai{Hai: h, isOpen: true}
	}
	if c.Tsumohai() != nil {
		hais[len(c.Tehai().Hais())+1] = &boardHai{Hai: c.Tsumohai(), isOpen: true}
	}

	head := 20
	// chii
	for _, meld := range c.Naki().Chiis() {
		for _, h := range meld {
			head--
			hais[head] = &boardHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// pon
	for _, meld := range c.Naki().Pons() {
		for _, h := range meld {
			head--
			hais[head] = &boardHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// ankan
	for _, meld := range c.Naki().AnKans() {
		for i, h := range meld {
			isOpen := true
			if i == 0 || i == 3 {
				isOpen = false
			}
			head--
			hais[head] = &boardHai{Hai: h, isOpen: isOpen, isDown: true}
		}
	}
	// minkan
	for _, meld := range c.Naki().MinKans() {
		for _, h := range meld {
			head--
			hais[head] = &boardHai{Hai: h, isOpen: true, isDown: true}
		}
	}

	// riichi
	if c.IsRiichi() {
		for _, h := range hais {
			if h != nil {
				h.isRiichi = true
			}
		}
	}
	return hais
}

func drawTehai(hais [20]*boardHai) string {
	strs := []string{"", "", "", ""}
	for _, h := range hais {

		if h == nil {
			strs[0] += "    "
			strs[1] += "    "
			strs[2] += "    "
			strs[3] += "    "
			continue
		}

		lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
		if h != nil && h.isRiichi {
			lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
		}
		if h.isDown {
			strs[0] += lines[0] + lines[1] + lines[1] + lines[2]
			strs[1] += lines[3] + h.Name() + lines[5]
			strs[2] += lines[6] + lines[7] + lines[7] + lines[8]
			strs[3] += lines[6] + lines[7] + lines[7] + lines[8]
		} else {
			if h.isOpen {
				strs[0] += "    "
				strs[1] += lines[0] + lines[1] + lines[1] + lines[2]
				strs[2] += lines[3] + h.Name() + lines[5]
				strs[3] += lines[6] + lines[7] + lines[7] + lines[8]

			} else {
				strs[0] += "    "
				strs[1] += lines[0] + lines[1] + lines[1] + lines[2]
				strs[2] += lines[3] + lines[4] + lines[4] + lines[5]
				strs[3] += lines[6] + lines[7] + lines[7] + lines[8]
			}
		}
	}
	return strings.Join(strs, "\n") + "\n"
}

func (t *boardImpl) draftHo(c player.Player) [20][20]*boardHai {
	hoHais := [20][20]*boardHai{}
	idx, err := t.MyTurn(c)
	if err != nil {
		panic(err)
	}
	myself := t.players[idx]
	shimoplayer := t.players[(idx+1)%t.maxNumberOfUser]
	toimen := t.players[(idx+2)%t.maxNumberOfUser]
	kamiplayer := t.players[(idx+3)%t.maxNumberOfUser]

	for i, h := range myself.Kawa().Hais() {
		hoHais[13+i/6][7+i%6] = &boardHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range shimoplayer.Kawa().Hais() {
		hoHais[12-i%6][13+i/6] = &boardHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range toimen.Kawa().Hais() {
		hoHais[6-i/6][12-i%6] = &boardHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range kamiplayer.Kawa().Hais() {
		hoHais[7+i%6][6-i/6] = &boardHai{Hai: h, isOpen: true, isDown: true}
	}

	return hoHais
}

func drawHo(hais [20][20]*boardHai) string {
	str := ""
	for i, _ := range hais {
		body := ""
		bottom := ""
		top := ""

		for j, h := range hais[i] {
			lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
			if (h != nil && h.isRiichi) ||
				(i != 0 && hais[i-1][j] != nil && hais[i-1][j].isRiichi) ||
				(i != len(hais)-1 && hais[i+1][j] != nil && hais[i+1][j].isRiichi) {
				lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
			}
			if i == 0 {
				if h != nil {
					top += lines[0] + lines[1] + lines[1] + lines[2]
				} else {
					top += "    "
				}
			}
			if h == nil {
				if i != 0 && hais[i-1][j] != nil {
					body += lines[6] + lines[7] + lines[7] + lines[8]
				} else {
					body += "    "
				}
				if i != len(hais)-1 && hais[i+1][j] != nil {
					bottom += lines[0] + lines[1] + lines[1] + lines[2]
				} else {
					bottom += "    "
				}
				continue
			}
			if h.isOpen && h.isDown {
				body += lines[3] + h.Name() + lines[5]
			} else {
				body += lines[3] + lines[4] + lines[4] + lines[5]
			}
			bottom += lines[6] + lines[7] + lines[7] + lines[8]
		}

		if i == len(hais)-1 {
			bottom += "\n"
			for _, h := range hais[i] {
				lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
				if h != nil && h.isRiichi {
					lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
				}
				if h == nil {
					bottom += "    "
				} else {
					bottom += lines[6] + lines[7] + lines[7] + lines[8]
				}
			}
		}

		if top != "" {
			str += top + "\n"
		}
		str += body + "\n"
		str += bottom + "\n"
	}
	return str
}
