package tehai

import (
	"mahjong/model/hai"
	"mahjong/model/hai/attribute"
)

func Kotsu(inHais []*hai.Hai) [][]*hai.Hai {
	cnt := map[*hai.Hai]int{}
	for _, h := range inHais {
		cnt[h]++
	}
	outHais := [][]*hai.Hai{}
	for k, v := range cnt {
		if v >= 3 {
			outHais = append(outHais, []*hai.Hai{k, k, k})
		}
	}

	return outHais
}

func Shuntsu(inHais []*hai.Hai) ([][]*hai.Hai, error) {
	outHais := [][]*hai.Hai{}
	for _, h := range inHais {
		if h.HasAttribute(&attribute.Jihai) {
			continue
		}
		suit, err := hai.HaitoSuits(h)
		if err != nil {
			return outHais, err
		}
		num, err := hai.HaitoI(h)
		if err != nil {
			return outHais, err
		}
		if num <= 7 && contain(inHais, suit[num]) && contain(inHais, suit[num+1]) {
			outHais = append(outHais, []*hai.Hai{h, suit[num], suit[num+1]})
		}
	}

	return outHais, nil
}

func contain(hais []*hai.Hai, hai *hai.Hai) bool {
	for _, h := range hais {
		if h == hai {
			return true
		}
	}
	return false
}
