package huro

import (
	"mahjong/hai"
)

type Huro interface {
	AddSet([]*hai.Hai) error
	AddHaiToSet(inTile *hai.Hai) error
}

type huroImpl struct {
	sets [][]*hai.Hai
}

func New() Huro {
	return &huroImpl{sets: [][]*hai.Hai{}}

}

func (h *huroImpl) AddSet(set []*hai.Hai) error {
	h.sets = append(h.sets, set)
	return nil
}

func (h *huroImpl) AddHaiToSet(inHai *hai.Hai) error {
	sidx := -1
	for idx, set := range h.sets {
		if len(set) != 3 {
			continue
		}
		if set[0] == inHai && set[1] == inHai && set[2] == inHai {
			sidx = idx
			break
		}
	}
	if sidx == -1 {
		return HuroNoSetFoundErr
	}
	h.sets[sidx] = append(h.sets[sidx], inHai)

	return nil
}
