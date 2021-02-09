package tehai

import (
	"mahjong/model/attribute"
	"mahjong/model/hai"
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
	FindChiiPairs(*hai.Hai) ([][2]*hai.Hai, error)
	FindPonPairs(*hai.Hai) ([][2]*hai.Hai, error)
	FindKanPairs(*hai.Hai) ([][3]*hai.Hai, error)
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

func (t *tehaiImpl) FindPonPairs(inHai *hai.Hai) ([][2]*hai.Hai, error) {
	pairs := [][2]*hai.Hai{}
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

func (t *tehaiImpl) FindKanPairs(inHai *hai.Hai) ([][3]*hai.Hai, error) {
	pairs := [][3]*hai.Hai{}
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

func (t *tehaiImpl) FindChiiPairs(inHai *hai.Hai) ([][2]*hai.Hai, error) {
	pairs := [][2]*hai.Hai{}
	if !inHai.HasAttribute(&attribute.Suhai) {
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

	// find right pair
	if num >= 3 && t.hasHai(suit[num-3]) && t.hasHai(suit[num-2]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-3], suit[num-2]})
	}
	// find center pair
	if num >= 2 && num <= 8 && t.hasHai(suit[num-2]) && t.hasHai(suit[num]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-2], suit[num]})
	}
	// find left pair
	if num <= 7 && t.hasHai(suit[num]) && t.hasHai(suit[num+1]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num], suit[num+1]})
	}

	return pairs, nil
}

func (t *tehaiImpl) hasHai(inHai *hai.Hai) bool {
	for _, h := range t.hais {
		if h == inHai {
			return true
		}
	}

	return false
}
