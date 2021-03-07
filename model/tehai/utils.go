package tehai

import (
	"mahjong/model/hai"
	"mahjong/model/hai/attribute"
	"sort"
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
	Sort(inHais)
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

func Sort(inHais []*hai.Hai) {
	sort.Slice(inHais, func(i int, j int) bool {
		return inHais[i].Name() < inHais[j].Name()
	})
}

func Unique(inHais []*hai.Hai) []*hai.Hai {
	uniq := map[*hai.Hai]bool{}
	for _, h := range inHais {
		uniq[h] = true
	}

	outHais := []*hai.Hai{}
	for k, _ := range uniq {
		outHais = append(outHais, k)
	}
	return outHais
}

func contain(hais []*hai.Hai, hai *hai.Hai) bool {
	for _, h := range hais {
		if h == hai {
			return true
		}
	}
	return false
}

func Machihai(a *hai.Hai, b *hai.Hai) ([]*hai.Hai, error) {
	outHais := []*hai.Hai{}
	if a == b {
		outHais = append(outHais, a)
	}

	// check suit
	if a.HasAttribute(&attribute.Manzu) && b.HasAttribute(&attribute.Manzu) ||
		a.HasAttribute(&attribute.Pinzu) && b.HasAttribute(&attribute.Pinzu) ||
		a.HasAttribute(&attribute.Souzu) && b.HasAttribute(&attribute.Souzu) {

		num1, err := hai.HaitoI(a)
		if err != nil {
			return outHais, err
		}
		num2, err := hai.HaitoI(b)
		if err != nil {
			return outHais, err
		}

		suit, err := hai.HaitoSuits(a)
		if err != nil {
			return outHais, err
		}
		minv := min(num1, num2)
		maxv := max(num1, num2)
		diff := maxv - minv
		switch diff {
		case 1:
			{
				if minv != 1 {
					outHais = append(outHais, suit[minv-2])
				}
				if maxv != 9 {
					outHais = append(outHais, suit[maxv])
				}

			}
		case 2:
			{
				outHais = append(outHais, suit[minv])
			}
		default:
			{

			}
		}
	}

	return outHais, nil
}

func min(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}
