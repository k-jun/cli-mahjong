package taku

import (
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/yama"
	"strings"
	"sync"
)

var (
	MaxNumberOfUsers = 4
)

type Taku interface {
	// setter
	SetWinIndex(int) error

	// game
	JoinCha(cha.Cha) (chan Taku, error)
	LeaveCha(cha.Cha) error
	Broadcast()

	// turn
	CurrentTurn() int
	NextTurn() int
	MyTurn(cha.Cha) (int, error)
	TurnEnd() error

	// last ho
	LastHo() (*hai.Hai, error)

	// action counter
	ActionCounter() int
	CancelAction(c cha.Cha) error
	TakeAction(cha.Cha, func(*hai.Hai) error) error

	// draw
	Draw(cha.Cha) string
}

func New(maxNOU int, y yama.Yama) Taku {
	return &takuImpl{
		chas:            []*takuCha{},
		yama:            y,
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		actionChas:      []*takuCha{},
		winIndex:        -1,
	}
}

type takuImpl struct {
	sync.Mutex
	chas            []*takuCha
	yama            yama.Yama
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
	actionChas      []*takuCha

	// win
	winIndex int
}

type takuCha struct {
	channel chan Taku
	cha.Cha
}

func (t *takuImpl) SetWinIndex(idx int) error {
	if idx >= len(t.chas) || idx < 0 {
		return TakuIndexOutOfRangeErr
	}
	t.winIndex = idx
	return nil
}

func (t *takuImpl) JoinCha(c cha.Cha) (chan Taku, error) {
	t.Lock()
	defer t.Unlock()
	if len(t.chas) >= t.maxNumberOfUser {
		return nil, TakuMaxNOUErr
	}

	if err := c.SetYama(t.yama); err != nil {
		return nil, err
	}
	channel := make(chan Taku, t.maxNumberOfUser*3)
	t.chas = append(t.chas, &takuCha{Cha: c, channel: channel})

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
	// tehai assign
	for _, tc := range t.chas {
		if err := tc.Haipai(); err != nil {
			return err
		}
	}

	// tsumo
	return t.chas[t.CurrentTurn()].Cha.Tsumo()
}

func (t *takuImpl) CurrentTurn() int {
	return t.turnIndex
}

func (t *takuImpl) MyTurn(c cha.Cha) (int, error) {
	for i, tc := range t.chas {
		if tc.Cha == c {
			return i, nil
		}
	}
	return -1, TakuChaNotFoundErr
}

func (t *takuImpl) NextTurn() int {
	return (t.turnIndex + 1) % t.maxNumberOfUser
}

func (t *takuImpl) TurnEnd() error {
	t.Lock()
	defer t.Unlock()
	err := t.setActionCounter()
	if err != nil {
		return err
	}

	if len(t.actionChas) == 0 {
		if err := t.turnChange(t.NextTurn()); err != nil {
			return err
		}
		if err := t.chas[t.CurrentTurn()].Cha.Tsumo(); err != nil {
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
	chas := []*takuCha{}

	inHai, err := t.chas[t.CurrentTurn()].Ho().Last()
	if err != nil {
		return err
	}
	for i, tc := range t.chas {
		if tc == t.chas[t.CurrentTurn()] {
			continue
		}

		actionCounter := 0
		if i == t.NextTurn() {
			pairs, err := tc.Cha.Tehai().FindChiiPairs(inHai)
			if err != nil {
				return err
			}
			actionCounter += len(pairs)
		}
		pairs, err := tc.Cha.Tehai().FindPonPairs(inHai)
		if err != nil {
			return err
		}
		actionCounter += len(pairs)
		kpairs, err := tc.Cha.Tehai().FindKanPairs(inHai)
		if err != nil {
			return err
		}
		actionCounter += len(kpairs)
		ok, err := tc.Cha.CanRon(inHai)
		if err != nil {
			return err
		}
		if ok {
			actionCounter += 1
		}

		if actionCounter != 0 {
			chas = append(chas, tc)
		}
	}
	t.actionChas = chas
	return nil
}

func (t *takuImpl) LastHo() (*hai.Hai, error) {
	return t.chas[t.CurrentTurn()].Ho().Last()
}

func (t *takuImpl) ActionCounter() int {
	return len(t.actionChas)
}

func (t *takuImpl) CancelAction(c cha.Cha) error {
	t.Lock()
	defer t.Unlock()
	if len(t.actionChas) == 0 {
		return nil
	}

	found := false
	for i, tc := range t.actionChas {
		if tc.Cha == c {
			found = true
			t.actionChas = append(t.actionChas[:i], t.actionChas[i+1:]...)
		}
	}
	if !found {
		return TakuChaNotFoundErr
	}

	if len(t.actionChas) == 0 {
		if err := t.turnChange(t.NextTurn()); err != nil {
			return err
		}
		if err := t.chas[t.CurrentTurn()].Cha.Tsumo(); err != nil {
			return err
		}
		go t.Broadcast()
	}
	return nil
}

func (t *takuImpl) TakeAction(c cha.Cha, action func(*hai.Hai) error) error {
	t.Lock()
	defer t.Unlock()
	if len(t.actionChas) == 0 {
		return TakuActionAlreadyTokenErr
	}

	found := false
	for _, tc := range t.actionChas {
		if tc.Cha == c {
			found = true
		}
	}
	if !found {
		return TakuChaNotFoundErr
	}

	t.actionChas = []*takuCha{}
	h, err := t.chas[t.CurrentTurn()].Ho().RemoveLast()
	if err != nil {
		return err
	}

	if err := action(h); err != nil {
		return err
	}

	turnIdx, _ := t.MyTurn(c)
	if err := t.turnChange(turnIdx); err != nil {
		return err
	}
	go t.Broadcast()
	return nil
}

type takuHai struct {
	*hai.Hai
	isOpen bool
	isDown bool
}

func (t *takuImpl) Draw(c cha.Cha) string {

	str := ""
	if t.winIndex != -1 {
		draftTehai := t.draftTehai(t.chas[t.winIndex])
		str += drawTehai(draftTehai)
		str += "\n GAME SET!!"
		return str
	}
	// tehais
	draftTehais := t.draftTehaiAll(c)
	str += drawTehai(draftTehais["toimen"])

	// ho
	draftHo := t.draftHo(c)

	for i, h := range draftTehais["kamicha"] {
		if h == nil {
			continue
		}
		draftHo[i][0] = h
	}
	for i, h := range draftTehais["shimocha"] {
		if h == nil {
			continue
		}
		draftHo[len(draftHo)-1-i][len(draftHo[i])-1] = h
	}
	str += drawHo(draftHo)

	// tehai
	str += drawTehai(draftTehais["jicha"])

	return str
}

func (t *takuImpl) draftTehaiAll(c cha.Cha) map[string][20]*takuHai {
	tehaiMap := map[string][20]*takuHai{}
	idx, err := t.MyTurn(c)
	if err != nil {
		panic(err)
	}
	jicha := t.chas[idx]
	shimocha := t.chas[(idx+1)%t.maxNumberOfUser]
	toimen := t.chas[(idx+2)%t.maxNumberOfUser]
	kamicha := t.chas[(idx+3)%t.maxNumberOfUser]
	tehaiMap["jicha"] = t.draftTehai(jicha)
	tehaiMap["shimocha"] = hideTehai(t.draftTehai(shimocha))
	tehaiMap["toimen"] = hideTehai(reverse(t.draftTehai(toimen)))
	tehaiMap["kamicha"] = hideTehai(t.draftTehai(kamicha))
	return tehaiMap
}

func reverse(s [20]*takuHai) [20]*takuHai {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func hideTehai(v [20]*takuHai) [20]*takuHai {
	hais := [20]*takuHai{}
	for i, h := range v {
		if h != nil && !h.isDown {
			h.isOpen = false
		}
		hais[i] = h
	}

	return hais
}

func (t *takuImpl) draftTehai(c cha.Cha) [20]*takuHai {
	hais := [20]*takuHai{}
	for i, h := range c.Tehai().Hais() {
		hais[i] = &takuHai{Hai: h, isOpen: true}
	}
	if c.Tsumohai() != nil {
		hais[len(c.Tehai().Hais())+1] = &takuHai{Hai: c.Tsumohai(), isOpen: true}
	}

	head := 20
	// chii
	for _, meld := range c.Huro().Chiis() {
		for _, h := range meld {
			head--
			hais[head] = &takuHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// pon
	for _, meld := range c.Huro().Pons() {
		for _, h := range meld {
			head--
			hais[head] = &takuHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// ankan
	for _, meld := range c.Huro().AnKans() {
		for i, h := range meld {
			isOpen := true
			if i == 0 || i == 3 {
				isOpen = false
			}
			head--
			hais[head] = &takuHai{Hai: h, isOpen: isOpen, isDown: true}
		}
	}
	// minkan
	for _, meld := range c.Huro().MinKans() {
		for _, h := range meld {
			head--
			hais[head] = &takuHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	return hais
}

func drawTehai(hais [20]*takuHai) string {
	strs := []string{"", "", "", ""}
	for _, h := range hais {
		if h == nil {
			strs[0] += "    "
			strs[1] += "    "
			strs[2] += "    "
			strs[3] += "    "
			continue
		}
		if h.isDown {
			strs[0] += "┌──┐"
			strs[1] += "│" + h.Name() + "│"
			strs[2] += "└──┘"
			strs[3] += "└──┘"
		} else {
			if h.isOpen {
				strs[0] += "    "
				strs[1] += "┌──┐"
				strs[2] += "│" + h.Name() + "│"
				strs[3] += "└──┘"

			} else {
				strs[0] += "    "
				strs[1] += "┌──┐"
				strs[2] += "│  │"
				strs[3] += "└──┘"
			}
		}
	}
	return strings.Join(strs, "\n") + "\n"
}

func (t *takuImpl) draftHo(c cha.Cha) [20][20]*takuHai {
	hoHais := [20][20]*takuHai{}
	idx, err := t.MyTurn(c)
	if err != nil {
		panic(err)
	}
	myself := t.chas[idx]
	shimocha := t.chas[(idx+1)%t.maxNumberOfUser]
	toimen := t.chas[(idx+2)%t.maxNumberOfUser]
	kamicha := t.chas[(idx+3)%t.maxNumberOfUser]

	for i, h := range myself.Ho().Hais() {
		hoHais[13+i/6][7+i%6] = &takuHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range shimocha.Ho().Hais() {
		hoHais[12-i%6][13+i/6] = &takuHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range toimen.Ho().Hais() {
		hoHais[6-i/6][12-i%6] = &takuHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range kamicha.Ho().Hais() {
		hoHais[7+i%6][6-i/6] = &takuHai{Hai: h, isOpen: true, isDown: true}
	}
	return hoHais
}

func drawHo(hais [20][20]*takuHai) string {
	str := ""
	for i, _ := range hais {
		body := ""
		bottom := ""
		top := ""
		if i == 0 {
			for _, h := range hais[i] {
				if h != nil {
					top += "┌──┐"
				} else {
					top += "    "
				}
			}
		}

		for j, h := range hais[i] {
			if h == nil {
				if i != 0 && hais[i-1][j] != nil {
					body += "└──┘"
				} else {
					body += "    "
				}
				if i != len(hais)-1 && hais[i+1][j] != nil {
					bottom += "┌──┐"
				} else {
					bottom += "    "
				}
				continue
			}
			if h.isOpen && h.isDown {
				body += "│" + h.Name() + "│"
			} else {
				body += "│  │"
			}
			bottom += "└──┘"
		}

		if i == len(hais)-1 {
			bottom += "\n"
			for _, h := range hais[i] {
				if h == nil {
					bottom += "    "
				} else {
					bottom += "└──┘"
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
