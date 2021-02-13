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
	// game
	JoinCha(cha.Cha) (chan Taku, error)
	LeaveCha(cha.Cha) error
	Broadcast()

	// turn
	CurrentTurn() int
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

func New(maxNOU int) Taku {
	return &takuImpl{
		chas:            []*takuCha{},
		turnIndex:       0,
		maxNumberOfUser: maxNOU,
		isPlaying:       true,
		actionChas:      []*takuCha{},
	}
}

type takuImpl struct {
	sync.Mutex
	chas            []*takuCha
	turnIndex       int
	maxNumberOfUser int
	isPlaying       bool
	actionChas      []*takuCha
}

type takuCha struct {
	channel chan Taku
	cha.Cha
}

func (t *takuImpl) JoinCha(c cha.Cha) (chan Taku, error) {
	t.Lock()
	defer t.Unlock()
	if len(t.chas) >= t.maxNumberOfUser {
		return nil, TakuMaxNOUErr
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
	// create yama
	y := yama.New()
	if err := y.Kan(); err != nil {
		return err
	}

	// tehai assign
	for _, tc := range t.chas {
		if err := tc.SetYama(y); err != nil {
			return err
		}
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

	if len(t.actionChas) == 0 {
		if err := t.turnChange(t.nextTurn()); err != nil {
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
		if i == t.nextTurn() {
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

		if err != nil {
			return err
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
		if err := t.turnChange(t.nextTurn()); err != nil {
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
	// tehais
	draftTehais := t.draftTehaiAll(c)
	str += drawTehai(draftTehais["toimen"])

	// ho
	draftHo := t.draftHo(c)
	str += drawHo(draftHo[:])

	// tehai
	str += drawTehai(draftTehais["jicha"])

	return str
}

func (t *takuImpl) draftTehaiAll(c cha.Cha) map[string][]*takuHai {
	tehaiMap := map[string][]*takuHai{}
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

func reverse(s []*takuHai) []*takuHai {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func hideTehai(v []*takuHai) []*takuHai {

	hais := []*takuHai{}
	for _, h := range v {
		if h != nil && !h.isDown {
			h.isOpen = false
		}
		hais = append(hais, h)
	}

	return hais
}

func (t *takuImpl) draftTehai(c cha.Cha) []*takuHai {
	hais := []*takuHai{}
	for _, h := range c.Tehai().Hais() {
		hais = append(hais, &takuHai{Hai: h, isOpen: true})
	}
	if c.Tsumohai() != nil {
		hais = append(hais, nil)
		hais = append(hais, &takuHai{Hai: c.Tsumohai(), isOpen: true})
	}
	if len(c.Huro().Chiis()) != 0 || len(c.Huro().Pons()) != 0 || len(c.Huro().AnKans()) != 0 || len(c.Huro().AnKans()) != 0 {
		hais = append(hais, nil)
	}

	// chii
	for _, meld := range c.Huro().Chiis() {
		for _, h := range meld {
			hais = append(hais, &takuHai{Hai: h, isOpen: true, isDown: true})
		}
	}
	// pon
	for _, meld := range c.Huro().Pons() {
		for _, h := range meld {
			hais = append(hais, &takuHai{Hai: h, isOpen: true, isDown: true})
		}
	}
	// ankan
	for _, meld := range c.Huro().AnKans() {
		for i, h := range meld {
			isOpen := true
			if i == 0 || i == 3 {
				isOpen = false
			}
			hais = append(hais, &takuHai{Hai: h, isOpen: isOpen, isDown: true})
		}
	}
	for _, meld := range c.Huro().MinKans() {
		for _, h := range meld {
			hais = append(hais, &takuHai{Hai: h, isOpen: true, isDown: true})
		}
	}
	return hais
}

func drawTehai(hais []*takuHai) string {
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
	return strings.Join(strs, "\n")
}

func (t *takuImpl) draftHo(c cha.Cha) [12][12]*hai.Hai {
	hoHais := [12][12]*hai.Hai{}
	idx, err := t.MyTurn(c)
	if err != nil {
		panic(err)
	}
	myself := t.chas[idx]
	shimocha := t.chas[(idx+1)%t.maxNumberOfUser]
	toimen := t.chas[(idx+2)%t.maxNumberOfUser]
	kamicha := t.chas[(idx+3)%t.maxNumberOfUser]

	for i, h := range myself.Ho().Hais() {
		hoHais[9+i/6][3+i%6] = h
	}
	for i, h := range shimocha.Ho().Hais() {
		hoHais[8-i%6][9+i/6] = h
	}
	for i, h := range toimen.Ho().Hais() {
		hoHais[2-i/6][8-i%6] = h
	}
	for i, h := range kamicha.Ho().Hais() {
		hoHais[3+i%6][2-i/6] = h
	}
	return hoHais
}

func drawHo(hais [][12]*hai.Hai) string {
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
			body += "│" + h.Name() + "│"
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
