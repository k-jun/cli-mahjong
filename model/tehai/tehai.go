package tehai

import (
	"mahjong/model/hai"
	"mahjong/model/hai/attribute"
	"sort"
)

var (
	MaxHaisLen = 13
)

type Tehai interface {
	Hais() []*hai.Hai

	Add(*hai.Hai) error
	Adds([]*hai.Hai) error
	Remove(*hai.Hai) (*hai.Hai, error)
	Removes([]*hai.Hai) ([]*hai.Hai, error)
	Replace(*hai.Hai, *hai.Hai) (*hai.Hai, error)

	ChiiPairs(*hai.Hai) ([][2]*hai.Hai, error)
	PonPairs(*hai.Hai) ([][2]*hai.Hai, error)
	MinKanPairs(*hai.Hai) ([][3]*hai.Hai, error)
	AnKanPairs(*hai.Hai) ([][4]*hai.Hai, error)
	RiichiHais(*hai.Hai) ([]*hai.Hai, error)

	CanChii(*hai.Hai) (bool, error)
	CanPon(*hai.Hai) (bool, error)
	CanMinKan(*hai.Hai) (bool, error)
	CanAnKan(*hai.Hai) (bool, error)
	CanRiichi(*hai.Hai) (bool, error)
	CanRon(*hai.Hai) (bool, error)

	Sort() error
}

type tehaiImpl struct {
	hais []*hai.Hai
}

func New() Tehai {
	return &tehaiImpl{hais: []*hai.Hai{}}
}

func (t *tehaiImpl) Hais() []*hai.Hai {
	return t.hais
}

func (t *tehaiImpl) Sort() error {
	sort.Slice(t.hais, func(i int, j int) bool {
		return t.hais[i].Name() < t.hais[j].Name()
	})
	return nil
}

func (t *tehaiImpl) Add(inHai *hai.Hai) error {
	if inHai == nil {
		return TehaiHaiIsNilErr
	}
	if len(t.hais) >= MaxHaisLen {
		return TehaiReachMaxHaiErr
	}

	t.hais = append(t.hais, inHai)
	return nil
}

func (t *tehaiImpl) Adds(inHais []*hai.Hai) error {
	for _, hai := range inHais {
		if err := t.Add(hai); err != nil {
			return err
		}
	}

	return nil
}

func (t *tehaiImpl) Remove(outHai *hai.Hai) (*hai.Hai, error) {
	for idx, hai := range t.hais {
		if hai == outHai {
			outHai = t.hais[idx]
			t.hais = append(t.hais[:idx], t.hais[idx+1:]...)
			return outHai, nil
		}

	}
	return nil, TehaiHaiNotFoundErr
}

func (t *tehaiImpl) Removes(outHais []*hai.Hai) ([]*hai.Hai, error) {
	hais := []*hai.Hai{}
	for _, outHai := range outHais {
		outHai, err := t.Remove(outHai)
		if err != nil {
			return outHais, err
		}
		hais = append(hais, outHai)

	}
	return outHais, nil
}

func (t *tehaiImpl) Replace(inHai *hai.Hai, outHai *hai.Hai) (*hai.Hai, error) {
	for idx, hai := range t.hais {
		if hai == outHai {
			outHai = t.hais[idx]
			t.hais[idx] = inHai
			return outHai, nil
		}
	}
	return nil, TehaiHaiNotFoundErr
}

func (t *tehaiImpl) CanChii(inHai *hai.Hai) (bool, error) {
	hais, err := t.ChiiPairs(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) ChiiPairs(inHai *hai.Hai) ([][2]*hai.Hai, error) {
	pairs := [][2]*hai.Hai{}
	if inHai == nil || !inHai.HasAttribute(&attribute.Suhai) {
		return pairs, nil
	}
	// detect suit
	suit, err := hai.HaitoSuits(inHai)
	if err != nil {
		return pairs, nil
	}
	// detect hai's number
	num, err := hai.HaitoI(inHai)
	if err != nil {
		return pairs, nil
	}

	//  right pair
	if num >= 3 && t.hasHai(suit[num-3]) && t.hasHai(suit[num-2]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-3], suit[num-2]})
	}
	//  center pair
	if num >= 2 && num <= 8 && t.hasHai(suit[num-2]) && t.hasHai(suit[num]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-2], suit[num]})
	}
	//  left pair
	if num <= 7 && t.hasHai(suit[num]) && t.hasHai(suit[num+1]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num], suit[num+1]})
	}

	return pairs, nil
}

func (t *tehaiImpl) CanPon(inHai *hai.Hai) (bool, error) {
	hais, err := t.PonPairs(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) PonPairs(inHai *hai.Hai) ([][2]*hai.Hai, error) {
	pairs := [][2]*hai.Hai{}
	if inHai == nil {
		return pairs, nil
	}
	cnt := map[*hai.Hai]int{}
	for _, h := range t.hais {
		cnt[h] += 1
	}

	for k, v := range cnt {
		if v >= 2 && k == inHai {
			pairs = append(pairs, [2]*hai.Hai{k, k})
		}
	}

	return pairs, nil
}

func (t *tehaiImpl) CanMinKan(inHai *hai.Hai) (bool, error) {
	hais, err := t.MinKanPairs(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) MinKanPairs(inHai *hai.Hai) ([][3]*hai.Hai, error) {
	pairs := [][3]*hai.Hai{}
	if inHai == nil {
		return pairs, nil
	}
	cnt := map[*hai.Hai]int{}
	for _, h := range t.hais {
		cnt[h] += 1
	}

	for k, v := range cnt {
		if v >= 3 && k == inHai {
			pairs = append(pairs, [3]*hai.Hai{k, k, k})
		}
	}
	return pairs, nil
}

func (t *tehaiImpl) CanAnKan(inHai *hai.Hai) (bool, error) {
	hais, err := t.AnKanPairs(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) AnKanPairs(inHai *hai.Hai) ([][4]*hai.Hai, error) {
	pairs := [][4]*hai.Hai{}
	if inHai == nil {
		return pairs, nil
	}
	cnt := map[*hai.Hai]int{}
	cnt[inHai]++
	for _, h := range t.hais {
		cnt[h]++
	}

	for k, v := range cnt {
		if v >= 4 {
			pairs = append(pairs, [4]*hai.Hai{k, k, k, k})
		}
	}
	return pairs, nil
}

func (t *tehaiImpl) CanRon(inHai *hai.Hai) (bool, error) {
	if inHai == nil {
		return false, nil
	}
	tehai := tehaiImpl{append([]*hai.Hai{inHai}, t.hais...)}
	tehai.Sort()

	// count by hai
	haisMap := map[*hai.Hai]int{}
	for _, h := range tehai.hais {
		haisMap[h]++
	}
	for k, v := range haisMap {
		if v < 2 {
			continue
		}
		// deep copy
		tehai := tehaiImpl{append([]*hai.Hai{}, tehai.hais...)}
		if _, err := tehai.Removes([]*hai.Hai{k, k}); err != nil {
			return false, err
		}

		// kotsu
		pairs := Kotsu(tehai.hais)
		for i := 0; i < 1<<len(pairs); i++ {
			// deep copy
			tehai := tehaiImpl{append([]*hai.Hai{}, tehai.hais...)}
			// remove
			for j := 0; j < len(pairs); j++ {
				if i>>j&1 == 1 {
					if _, err := tehai.Removes(pairs[j]); err != nil {
						return false, err
					}
				}
			}

			// shuntsu
			pairs, err := Shuntsu(tehai.hais)
			if err != nil {
				return false, err
			}
			// remove
			for j := 0; j < len(pairs); j++ {
				if !(tehai.hasHai(pairs[j][0]) && tehai.hasHai(pairs[j][1]) && tehai.hasHai(pairs[j][2])) {
					continue
				}
				if _, err := tehai.Removes(pairs[j]); err != nil {
					return false, err
				}
			}

			if len(tehai.hais) == 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

func (t *tehaiImpl) CanRiichi(inHai *hai.Hai) (bool, error) {
	if inHai == nil {
		return false, nil
	}
	hais, err := t.RiichiHais(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) RiichiHais(inHai *hai.Hai) ([]*hai.Hai, error) {
	outHais := []*hai.Hai{}
	if inHai == nil {
		return outHais, nil
	}
	hais := append([]*hai.Hai{}, t.hais...)
	hais = append(hais, inHai)
	for i, outHai := range hais {
		// deep copy, and remove outHai
		hais_copy := append([]*hai.Hai{}, hais[:i]...)
		hais_copy = append(hais_copy, hais[i+1:]...)

		// check machihai
		tehai := tehaiImpl{hais: hais_copy}
		hais, err := tehai.Machihai()
		if err != nil {
			return outHais, err
		}
		if len(hais) != 0 {
			outHais = append(outHais, outHai)
		}
	}

	return outHais, nil
}

func (t *tehaiImpl) Machihai() ([]*hai.Hai, error) {
	machihai := []*hai.Hai{}

	for _, h := range hai.All {
		ok, err := t.CanRon(h)
		if err != nil {
			return machihai, err
		}
		if ok {
			machihai = append(machihai, h)
		}
	}
	return machihai, nil
}

func (t *tehaiImpl) hasHai(inHai *hai.Hai) bool {
	for _, h := range t.hais {
		if h == inHai {
			return true
		}
	}

	return false
}
