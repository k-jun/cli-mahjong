package tehai

import (
	"mahjong/attribute"
	"mahjong/hai"
)

var (
	MaxHaisLen = 13
)

type Tehai interface {
	Len() int
	Add(*hai.Hai) error
	Adds([]*hai.Hai) error
	Remove(*hai.Hai) (*hai.Hai, error)
	Removes([]*hai.Hai) ([]*hai.Hai, error)
	Replace(*hai.Hai, *hai.Hai) (*hai.Hai, error)
	FindChiPairs(*hai.Hai) [][2]*hai.Hai
	FindPonPairs(*hai.Hai) [][2]*hai.Hai
	FindKanPairs(*hai.Hai) [][3]*hai.Hai
	Hais() []*hai.Hai
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

func (t *tehaiImpl) Len() int {
	return len(t.hais)
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
			t.hais[idx] = t.hais[0]
			t.hais = t.hais[1:]
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

func (t *tehaiImpl) FindPonPairs(inHai *hai.Hai) [][2]*hai.Hai {
	pairs := [][2]*hai.Hai{}
	cnt := map[*hai.Hai]int{}
	for _, h := range t.hais {
		cnt[h] += 1
	}

	for k, v := range cnt {
		if v >= 2 {
			pairs = append(pairs, [2]*hai.Hai{k, k})
		}
	}

	return pairs
}

func (t *tehaiImpl) FindKanPairs(inHai *hai.Hai) [][3]*hai.Hai {
	pairs := [][3]*hai.Hai{}
	cnt := map[*hai.Hai]int{}
	for _, h := range t.hais {
		cnt[h] += 1
	}

	for k, v := range cnt {
		if v >= 3 {
			pairs = append(pairs, [3]*hai.Hai{k, k, k})
		}
	}

	return pairs
}

func (t *tehaiImpl) FindChiPairs(inHai *hai.Hai) [][2]*hai.Hai {
	pairs := [][2]*hai.Hai{}
	if !inHai.HasAttribute(&attribute.Suhai) {
		return pairs
	}
	// detect suit
	suit := []*hai.Hai{}
	if inHai.HasAttribute(&attribute.Manzu) {
		suit = hai.Manzu
	}
	if inHai.HasAttribute(&attribute.Pinzu) {
		suit = hai.Pinzu
	}
	if inHai.HasAttribute(&attribute.Souzu) {
		suit = hai.Souzu
	}
	// detect hai's number
	num := 0
	for idx, h := range suit {
		if h == inHai {
			num = idx
		}
	}

	// find right pair
	if num >= 3 && t.hasHai(suit[num-2]) && t.hasHai(suit[num-1]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-2], suit[num-1]})
	}
	// find center pair
	if num >= 2 && num <= 8 && t.hasHai(suit[num-1]) && t.hasHai(suit[num+1]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num-1], suit[num+1]})
	}
	// find left pair
	if num <= 7 && t.hasHai(suit[num+1]) && t.hasHai(suit[num+2]) {
		pairs = append(pairs, [2]*hai.Hai{suit[num+1], suit[num+2]})
	}

	return pairs
}

func (t *tehaiImpl) hasHai(inHai *hai.Hai) bool {
	for _, h := range t.hais {
		if h == inHai {
			return true
		}
	}

	return false
}
