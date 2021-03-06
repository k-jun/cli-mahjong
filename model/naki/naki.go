package naki

import (
	"mahjong/model/hai"
)

type Naki interface {
	SetPon([3]*hai.Hai) error
	SetChii([3]*hai.Hai) error
	SetMinKan([4]*hai.Hai) error
	SetAnKan([4]*hai.Hai) error

	Pons() [][3]*hai.Hai
	Chiis() [][3]*hai.Hai
	MinKans() [][4]*hai.Hai
	AnKans() [][4]*hai.Hai

	CanKakan(*hai.Hai) bool
	Kakan(*hai.Hai) error
}

type nakiImpl struct {
	pons    [][3]*hai.Hai
	chiis   [][3]*hai.Hai
	minKans [][4]*hai.Hai
	anKans  [][4]*hai.Hai
}

func New() Naki {
	return &nakiImpl{
		pons:    [][3]*hai.Hai{},
		chiis:   [][3]*hai.Hai{},
		minKans: [][4]*hai.Hai{},
		anKans:  [][4]*hai.Hai{},
	}
}

func (h *nakiImpl) Pons() [][3]*hai.Hai {
	return h.pons
}
func (h *nakiImpl) Chiis() [][3]*hai.Hai {
	return h.chiis
}
func (h *nakiImpl) MinKans() [][4]*hai.Hai {
	return h.minKans
}
func (h *nakiImpl) AnKans() [][4]*hai.Hai {
	return h.anKans
}

func (h *nakiImpl) SetPon(hais [3]*hai.Hai) error {
	h.pons = append(h.pons, hais)
	return nil
}

func (h *nakiImpl) SetChii(hais [3]*hai.Hai) error {
	h.chiis = append(h.chiis, hais)
	return nil
}

func (h *nakiImpl) SetMinKan(hais [4]*hai.Hai) error {
	h.minKans = append(h.minKans, hais)
	return nil
}
func (h *nakiImpl) SetAnKan(hais [4]*hai.Hai) error {
	h.anKans = append(h.anKans, hais)
	return nil
}

func (h *nakiImpl) CanKakan(inHai *hai.Hai) bool {
	for _, pon := range h.pons {
		if pon[0] == inHai {
			return true
		}
	}
	return false
}

func (h *nakiImpl) Kakan(inHai *hai.Hai) error {
	for idx, pon := range h.pons {
		if pon[0] == inHai {
			h.pons[idx] = h.pons[0]
			h.pons = h.pons[1:]

			set := [4]*hai.Hai{}
			set[0], set[1], set[2], set[3] = pon[0], pon[1], pon[2], inHai
			h.minKans = append(h.minKans, set)
			return nil
		}
	}

	return NakiNotFoundErr
}
