package player

import (
	"mahjong/model/hai"
	"mahjong/model/hai/attribute"
	"mahjong/model/kawa"
	"mahjong/model/naki"
	"mahjong/model/tehai"
	"mahjong/model/yama"
	"sort"

	"github.com/google/uuid"
)

type Player interface {
	// getter
	Tehai() tehai.Tehai
	Kawa() kawa.Kawa
	Tsumohai() *hai.Hai
	Naki() naki.Naki
	IsRiichi() bool
	// setter
	SetYama(yama.Yama) error
	// judger
	CanTsumoAgari() (bool, error)
	CanRon(*hai.Hai) (bool, error)
	FindRiichiHai() ([]*hai.Hai, error)
	FindNakiActions(*hai.Hai) ([]Action, error)

	Tsumo() error
	Dahai(*hai.Hai) error
	Haipai() error
	Chii(*hai.Hai, [2]*hai.Hai) error
	Pon(*hai.Hai, [2]*hai.Hai) error
	Kan(*hai.Hai, [3]*hai.Hai) error
	Kakan(*hai.Hai) error
	Riichi(*hai.Hai) error
}

type Action string

var (
	Chii       Action = "chii"
	Pon        Action = "pon"
	Kan        Action = "kan"
	Ron        Action = "ron"
	Riichi     Action = "riichi"
	Tsumo      Action = "tsumo"
	AllActions        = []Action{Chii, Pon, Kan, Ron, Riichi, Tsumo}
)

func AtoAction(s string) (Action, error) {
	for _, a := range AllActions {
		if string(a) == s {
			return a, nil
		}
	}

	return "", PlayerActionInvalidErr
}

type playerImpl struct {
	id       uuid.UUID
	tsumohai *hai.Hai
	kawa     kawa.Kawa
	tehai    tehai.Tehai
	naki     naki.Naki
	yama     yama.Yama
	isRiichi bool
}

func New(id uuid.UUID, k kawa.Kawa, t tehai.Tehai, n naki.Naki) Player {
	return &playerImpl{
		id:       id,
		kawa:     k,
		tehai:    t,
		naki:     n,
		yama:     nil,
		isRiichi: false,
	}
}

func (c *playerImpl) Tehai() tehai.Tehai {
	return c.tehai
}

func (c *playerImpl) Kawa() kawa.Kawa {
	return c.kawa
}

func (c *playerImpl) Naki() naki.Naki {
	return c.naki
}

func (c *playerImpl) Tsumohai() *hai.Hai {
	return c.tsumohai
}

func (c *playerImpl) IsRiichi() bool {
	return c.isRiichi
}

func (c *playerImpl) Tsumo() error {
	if c.tsumohai != nil {
		return PlayerAlreadyHaveTsumohaiErr
	}

	tsumohai, err := c.yama.Draw()
	if err != nil {
		return err
	}

	c.tsumohai = tsumohai
	return nil
}

func (c *playerImpl) Dahai(outHai *hai.Hai) error {
	var err error
	if c.isRiichi && outHai != c.tsumohai {
		return PlayerAlreadyRiichiErr
	}
	if outHai != c.tsumohai {
		if c.tsumohai == nil {
			outHai, err = c.tehai.Remove(outHai)
		} else {
			outHai, err = c.tehai.Replace(c.tsumohai, outHai)
		}

		if err != nil {
			return err
		}
		if err := c.tehai.Sort(); err != nil {
			return err
		}
	}
	c.tsumohai = nil

	return c.kawa.Add(outHai)
}

func (c *playerImpl) SetYama(y yama.Yama) error {
	if c.yama != nil {
		return PlayerAlreadyHaveYamaErr
	}
	c.yama = y
	return nil
}

func (c *playerImpl) Haipai() error {
	if len(c.tehai.Hais()) != 0 {
		return PlayerAlreadyDidHaipaiErr
	}

	for i := 0; i < tehai.MaxHaisLen; i++ {
		tsumoHai, err := c.yama.Draw()
		if err != nil {
			return err
		}
		if err := c.tehai.Add(tsumoHai); err != nil {
			return err
		}
	}

	if err := c.tehai.Sort(); err != nil {
		return err
	}
	return nil
}

func (c *playerImpl) Chii(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.naki.SetChii(meld)
}

func (c *playerImpl) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.naki.SetPon(meld)
}

func (c *playerImpl) Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1], outHais[2]})
	if err != nil {
		return err
	}
	meld := [4]*hai.Hai{inHai, hais[0], hais[1], hais[2]}

	return c.naki.SetMinKan(meld)
}

func (c *playerImpl) Kakan(inHai *hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	return c.naki.Kakan(inHai)
}

func (c *playerImpl) Riichi(inHai *hai.Hai) error {
	err := c.Dahai(inHai)
	if err != nil {
		return err
	}
	c.isRiichi = true
	return nil
}

func (c *playerImpl) FindRiichiHai() ([]*hai.Hai, error) {
	outHais := []*hai.Hai{}
	if len(c.naki.Chiis()) != 0 || len(c.naki.Pons()) != 0 || len(c.naki.MinKans()) != 0 || c.tsumohai == nil || c.isRiichi {
		return outHais, nil
	}

	hais := c.tehai.Hais()
	hais = append(hais, c.tsumohai)

	for _, eh := range hais {
		hs := append([]*hai.Hai{}, hais...)
		hs = removeHai(hs, eh)
		riichi, err := isTempan(hs)
		if err != nil {
			return outHais, err
		}
		if riichi && !haiContain(outHais, eh) {
			outHais = append(outHais, eh)
		}
	}
	return outHais, nil
}

func (c *playerImpl) FindAnKanHai() ([][3]*hai.Hai, error) {
	return c.tehai.FindKanPairs(c.tsumohai)
}

func (c *playerImpl) CanTsumoAgari() (bool, error) {
	if c.tsumohai == nil {
		return false, nil
	}
	hais := c.tehai.Hais()
	hais = append(hais, c.tsumohai)
	cnt := map[*hai.Hai]int{}
	for _, h := range hais {
		cnt[h] += 1
	}
	for k, v := range cnt {
		if v < 2 {
			continue
		}
		// deep copy
		hais := append([]*hai.Hai{}, hais...)
		hais = removeHais(hais, []*hai.Hai{k, k})
		hais, err := removeMentsus(hais)
		if err != nil {
			return false, err
		}
		if len(hais) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (c *playerImpl) CanRon(inHai *hai.Hai) (bool, error) {
	hais := c.tehai.Hais()
	hais = append(hais, inHai)
	cnt := map[*hai.Hai]int{}
	for _, h := range hais {
		cnt[h] += 1
	}
	for k, v := range cnt {
		if v < 2 {
			continue
		}
		// deep copy
		hais := append([]*hai.Hai{}, hais...)
		hais = removeHais(hais, []*hai.Hai{k, k})
		hais, err := removeMentsus(hais)
		if err != nil {
			return false, err
		}
		if len(hais) == 0 {
			return true, nil
		}
	}
	return false, nil

}

func (c *playerImpl) FindNakiActions(inHai *hai.Hai) ([]Action, error) {
	actions := []Action{}
	// chii
	pairs, err := c.tehai.FindChiiPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(pairs) != 0 {
		actions = append(actions, Chii)
	}

	// pon
	pairs, err = c.tehai.FindPonPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(pairs) != 0 {
		actions = append(actions, Pon)
	}

	// kan
	kanpairs, err := c.tehai.FindKanPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(kanpairs) != 0 {
		actions = append(actions, Kan)
	}

	// ron
	isRon, err := c.CanRon(inHai)
	if err != nil {
		return actions, err
	}
	if isRon {
		actions = append(actions, Ron)
	}
	return actions, nil
}

func haiContain(a []*hai.Hai, h *hai.Hai) bool {
	for _, hi := range a {
		if h == hi {
			return true
		}
	}
	return false

}

func isTempan(hais []*hai.Hai) (bool, error) {
	cnt := map[*hai.Hai]int{}
	for _, h := range hais {
		cnt[h] += 1
	}

	haisc := append([]*hai.Hai{}, hais...)
	haisc, err := removeMentsus(haisc)
	if err != nil {
		return false, err
	}
	if len(haisc) == 1 {
		return true, nil
	}

	for k, v := range cnt {
		hais := append([]*hai.Hai{}, hais...)
		if v < 2 {
			continue
		}
		// 両面, 嵌張, 双碰, 辺張
		hais = removeHais(hais, []*hai.Hai{k, k})
		hais, err := removeMentsus(hais)
		if err != nil {
			return false, err
		}
		if len(hais) == 2 && hasMati([2]*hai.Hai{hais[0], hais[1]}) {
			return true, nil
		}

	}

	return false, nil
}

func hasMati(hais [2]*hai.Hai) bool {
	if hais[0] == hais[1] {
		return true
	}
	// jihai
	if hais[0].HasAttribute(&attribute.Jihai) && hais[1].HasAttribute(&attribute.Jihai) {
		return hais[0] == hais[1]
	}

	// suhai
	num1, err := hai.HaitoI(hais[0])
	if err != nil {
		return false
	}
	num2, err := hai.HaitoI(hais[1])
	if err != nil {
		return false
	}
	if num2 > num1 {
		return num2-num1 <= 2
	}
	return num1-num2 <= 2
}

func removeHais(hais []*hai.Hai, outHais []*hai.Hai) []*hai.Hai {
	for _, hai := range outHais {
		hais = removeHai(hais, hai)
	}
	return hais
}

func removeHai(hais []*hai.Hai, hai *hai.Hai) []*hai.Hai {
	for i, h := range hais {
		if h == hai {
			hais = append(hais[:i], hais[i+1:]...)
			return hais
		}
	}
	panic(PlayerHaiNotFoundErr)
}

func removeMentsus(hais []*hai.Hai) ([]*hai.Hai, error) {
	for {
		kotsu, err := findKotsu(hais)
		if err != nil {
			return hais, err
		}
		if len(kotsu) != 0 {
			hais = removeHais(hais, kotsu)
		}
		shuntsu, err := findShuntsu(hais)
		if err != nil {
			return hais, err
		}
		if len(shuntsu) != 0 {
			hais = removeHais(hais, shuntsu)
		}

		if len(shuntsu) == 0 && len(kotsu) == 0 {
			break
		}
	}
	return hais, nil
}

func findShuntsu(hais []*hai.Hai) ([]*hai.Hai, error) {
	sort.Slice(hais, func(i int, j int) bool {
		return hais[i].Name() < hais[j].Name()
	})
	for _, h := range hais {
		if h.HasAttribute(&attribute.Jihai) {
			continue
		}
		suit, err := hai.HaitoSuits(h)
		if err != nil {
			return hais, err
		}
		num, err := hai.HaitoI(h)
		if err != nil {
			return hais, err
		}
		if num <= 7 && hasHai(hais, suit[num]) && hasHai(hais, suit[num+1]) {
			return []*hai.Hai{h, suit[num], suit[num+1]}, nil
		}
	}
	return []*hai.Hai{}, nil
}

func findKotsu(hais []*hai.Hai) ([]*hai.Hai, error) {
	cnt := map[*hai.Hai]int{}
	for _, h := range hais {
		cnt[h] += 1
		if cnt[h] >= 3 {
			return []*hai.Hai{h, h, h}, nil
		}
	}
	return []*hai.Hai{}, nil
}

func hasHai(hais []*hai.Hai, hai *hai.Hai) bool {
	for _, h := range hais {
		if h == hai {
			return true
		}
	}
	return false
}
