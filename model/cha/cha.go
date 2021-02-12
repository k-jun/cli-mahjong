package cha

import (
	"log"
	"mahjong/model/attribute"
	"mahjong/model/hai"
	"mahjong/model/ho"
	"mahjong/model/huro"
	"mahjong/model/tehai"
	"mahjong/model/yama"
	"sort"

	"github.com/google/uuid"
)

type Cha interface {
	// getter
	Tehai() tehai.Tehai
	Ho() ho.Ho
	Tsumohai() *hai.Hai
	Huro() huro.Huro
	// setter
	SetYama(yama.Yama) error
	// judger
	CanTsumoAgari() (bool, error)
	FindRiichiHai() ([]*hai.Hai, error)
	FindHuroActions(*hai.Hai) ([]huro.HuroAction, error)

	Tsumo() error
	Dahai(*hai.Hai) error
	Haipai() error
	Chii(*hai.Hai, [2]*hai.Hai) error
	Pon(*hai.Hai, [2]*hai.Hai) error
	Kan(*hai.Hai, [3]*hai.Hai) error
	Kakan(*hai.Hai) error
}

type chaImpl struct {
	id       uuid.UUID
	tsumohai *hai.Hai
	ho       ho.Ho
	tehai    tehai.Tehai
	huro     huro.Huro
	yama     yama.Yama
}

func New(id uuid.UUID, ho ho.Ho, t tehai.Tehai, y yama.Yama, hu huro.Huro) Cha {
	return &chaImpl{
		id:    id,
		ho:    ho,
		tehai: t,
		yama:  y,
		huro:  hu,
	}
}

func (c *chaImpl) Tehai() tehai.Tehai {
	return c.tehai
}

func (c *chaImpl) Ho() ho.Ho {
	return c.ho
}

func (c *chaImpl) Huro() huro.Huro {
	return c.huro
}

func (c *chaImpl) Tsumohai() *hai.Hai {
	return c.tsumohai
}

func (c *chaImpl) Tsumo() error {
	if c.tsumohai != nil {
		return ChaAlreadyHaveTsumohaiErr
	}

	tsumohai, err := c.yama.Draw()
	if err != nil {
		return err
	}

	c.tsumohai = tsumohai
	return nil
}

func (c *chaImpl) Dahai(outHai *hai.Hai) error {
	var err error
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

	return c.ho.Add(outHai)
}

func (c *chaImpl) SetYama(y yama.Yama) error {
	if c.yama != nil {
		return ChaAlreadyHaveYamaErr
	}
	c.yama = y
	return nil
}

func (c *chaImpl) Haipai() error {
	if len(c.tehai.Hais()) != 0 {
		return ChaAlreadyDidHaipaiErr
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

func (c *chaImpl) Chii(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.huro.SetChii(meld)
}

func (c *chaImpl) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.huro.SetPon(meld)
}

func (c *chaImpl) Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [4]*hai.Hai{inHai, hais[0], hais[1], hais[2]}

	return c.huro.SetMinKan(meld)
}

func (c *chaImpl) Kakan(inHai *hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	return c.huro.Kakan(inHai)
}

func (c *chaImpl) FindRiichiHai() ([]*hai.Hai, error) {
	outHais := []*hai.Hai{}
	if len(c.huro.Chiis()) != 0 || len(c.huro.Pons()) != 0 || len(c.huro.MinKans()) != 0 || c.tsumohai == nil {
		return outHais, nil
	}

	hais := c.tehai.Hais()
	hais = append(hais, c.tsumohai)

	for _, eh := range hais {
		hs := append([]*hai.Hai{}, hais...)
		hs = removeHai(hs, eh)
		riichi, err := isRiichi(hs)
		if err != nil {
			return outHais, err
		}
		if riichi && !haiContain(outHais, eh) {
			outHais = append(outHais, eh)
		}
	}
	return outHais, nil
}

func (c *chaImpl) CanTsumoAgari() (bool, error) {
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
		for {
			kotsu, err := findKotsu(hais)
			if err != nil {
				return false, err
			}
			shuntsu, err := findShuntsu(hais)
			if err != nil {
				return false, err
			}
			if len(kotsu) != 0 {
				hais = removeHais(hais, kotsu)
			}
			if len(shuntsu) != 0 {
				hais = removeHais(hais, shuntsu)
			}

			if len(shuntsu) == 0 && len(kotsu) == 0 {
				break
			}
		}
		if len(hais) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (c *chaImpl) FindHuroActions(inHai *hai.Hai) ([]huro.HuroAction, error) {
	actions := []huro.HuroAction{}

	// chii
	pairs, err := c.tehai.FindChiiPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(pairs) != 0 {
		actions = append(actions, huro.Chii)
	}

	// pon
	pairs, err = c.tehai.FindPonPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(pairs) != 0 {
		actions = append(actions, huro.Pon)
	}

	// kan
	kanpairs, err := c.tehai.FindKanPairs(inHai)
	if err != nil {
		return actions, err
	}
	if len(kanpairs) != 0 {
		actions = append(actions, huro.Kan)
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

func isRiichi(hais []*hai.Hai) (bool, error) {
	cnt := map[*hai.Hai]int{}
	for _, h := range hais {
		cnt[h] += 1
	}

	for k, v := range cnt {
		hais := append([]*hai.Hai{}, hais...)
		if v < 2 {
			// 単騎
			for {
				kotsu, err := findKotsu(hais)
				if err != nil {
					return false, err
				}
				shuntsu, err := findShuntsu(hais)
				if err != nil {
					return false, err

				}
				if len(kotsu) != 0 {
					hais = removeHais(hais, kotsu)
				}
				if len(shuntsu) != 0 {
					hais = removeHais(hais, shuntsu)
				}

				if len(shuntsu) == 0 && len(kotsu) == 0 {
					break
				}
			}
			if len(hais) == 1 {
				return true, nil
			}
		} else {
			// 両面, 嵌張, 双碰, 辺張
			hais = removeHais(hais, []*hai.Hai{k, k})
			for {
				kotsu, err := findKotsu(hais)
				if err != nil {
					return false, err
				}
				syuntu, err := findShuntsu(hais)
				if err != nil {
					return false, err
				}
				if len(kotsu) != 0 {
					hais = removeHais(hais, kotsu)
				}
				if len(syuntu) != 0 {
					hais = removeHais(hais, syuntu)
				}

				if len(syuntu) == 0 && len(kotsu) == 0 {
					break
				}
			}
			if len(hais) == 2 && hasMati([2]*hai.Hai{hais[0], hais[1]}) {
				return true, nil
			}
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
	for _, h := range hais {
		log.Println("h", h)
	}
	log.Println("outHais", outHais)
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
	log.Print("hais:")
	for _, h := range hais {
		log.Print(h.Name())
	}
	log.Println("")
	log.Println("hai:", hai)
	panic(ChaHaiNotFoundErr)
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
