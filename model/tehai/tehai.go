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
	KanPairs(*hai.Hai) ([][3]*hai.Hai, error)
	RiichiHais(*hai.Hai) ([]*hai.Hai, error)

	CanChii(*hai.Hai) (bool, error)
	CanPon(*hai.Hai) (bool, error)
	CanKan(*hai.Hai) (bool, error)
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

func (t *tehaiImpl) CanKan(inHai *hai.Hai) (bool, error) {
	hais, err := t.KanPairs(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) KanPairs(inHai *hai.Hai) ([][3]*hai.Hai, error) {
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

func (t *tehaiImpl) CanRon(inHai *hai.Hai) (bool, error) {
	machihai, err := t.Machihai()
	if err != nil {
		return false, err
	}
	for _, h := range machihai {
		if h == inHai {
			return true, nil
		}
	}
	return false, nil
}

func (t *tehaiImpl) CanRiichi(inHai *hai.Hai) (bool, error) {

	hais, err := t.RiichiHais(inHai)
	return len(hais) != 0, err
}

func (t *tehaiImpl) RiichiHais(inHai *hai.Hai) ([]*hai.Hai, error) {
	outHais := []*hai.Hai{}
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
	// TODO 七対子 国士無双

	// count by hai
	hais_map := map[*hai.Hai]int{}
	for _, h := range t.hais {
		hais_map[h]++
	}

	for k, v := range hais_map {
		// deep copy
		hais := append([]*hai.Hai{}, t.hais...)
		tehai := tehaiImpl{hais}
		if v >= 2 {
			// head remove
			if _, err := tehai.Removes([]*hai.Hai{k, k}); err != nil {
				return machihai, err
			}
		}
		// kotsu
		pairs := Kotsu(tehai.hais)
		for i := 0; i < 1<<len(pairs); i++ {
			// deep copy
			hais := append([]*hai.Hai{}, tehai.hais...)
			tehaiKotsu := tehaiImpl{hais}
			// remove
			for j := 0; j < len(pairs); j++ {
				if i>>j&1 == 1 {
					if _, err := tehaiKotsu.Removes(pairs[j]); err != nil {
						return machihai, err
					}
				}
			}

			// shuntsu
			print(tehaiKotsu.hais)
			forward := Sort(append([]*hai.Hai{}, tehaiKotsu.hais...))
			backward := ReverseSort(append([]*hai.Hai{}, tehaiKotsu.hais...))
			sorts := [][]*hai.Hai{forward, backward}

			for _, sortedHais := range sorts {
				tehaiShuntsu := tehaiImpl{sortedHais}
				pairs, err := Shuntsu(sortedHais)
				print(tehaiShuntsu.hais)
				if err != nil {
					return machihai, err
				}
				// remove
				for j := 0; j < len(pairs); j++ {
					if !(tehaiShuntsu.hasHai(pairs[j][0]) && tehaiShuntsu.hasHai(pairs[j][1]) && tehaiShuntsu.hasHai(pairs[j][2])) {
						continue
					}
					if _, err := tehaiShuntsu.Removes(pairs[j]); err != nil {
						// skip if not exist
						if err == TehaiHaiNotFoundErr {
							continue
						}
						return machihai, err
					}
				}

				// machi
				// tanki
				if len(tehaiShuntsu.hais) == 1 {
					machihai = append(machihai, tehaiShuntsu.hais[0])
				}
				// ryanmen etc..
				if len(tehaiShuntsu.hais) == 2 {
					hais, err := Machihai(tehaiShuntsu.hais[0], tehaiShuntsu.hais[1])
					if err != nil {
						return machihai, err
					}
					for _, h := range hais {
						machihai = append(machihai, h)
					}
				}
			}
		}
	}

	machihai = Unique(machihai)
	machihai = Sort(machihai)
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
