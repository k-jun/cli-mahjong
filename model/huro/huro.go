package huro

import (
	"mahjong/model/hai"
)

type HuroAction string

var (
	Chii HuroAction = "Chii"
	Pon  HuroAction = "Pon"
	Kan  HuroAction = "Kan"
	Ron  HuroAction = "Ron"
)

type Huro interface {
	SetPon([3]*hai.Hai) error
	SetChii([3]*hai.Hai) error
	SetMinKan([4]*hai.Hai) error
	SetAnKan([4]*hai.Hai) error
	Kakan(*hai.Hai) error
	Pons() [][3]*hai.Hai
	Chiis() [][3]*hai.Hai
	MinKans() [][4]*hai.Hai
	AnKans() [][4]*hai.Hai
}

type huroImpl struct {
	pons    [][3]*hai.Hai
	chiis   [][3]*hai.Hai
	minKans [][4]*hai.Hai
	ankans  [][4]*hai.Hai
}

func New() Huro {
	return &huroImpl{
		pons:    [][3]*hai.Hai{},
		chiis:   [][3]*hai.Hai{},
		minKans: [][4]*hai.Hai{},
		ankans:  [][4]*hai.Hai{},
	}
}

func (h *huroImpl) Pons() [][3]*hai.Hai {
	return h.pons
}
func (h *huroImpl) Chiis() [][3]*hai.Hai {
	return h.chiis
}
func (h *huroImpl) MinKans() [][4]*hai.Hai {
	return h.minKans
}
func (h *huroImpl) AnKans() [][4]*hai.Hai {
	return h.ankans
}

func (h *huroImpl) SetPon(hais [3]*hai.Hai) error {
	h.pons = append(h.pons, hais)
	return nil
}

func (h *huroImpl) SetChii(hais [3]*hai.Hai) error {
	h.chiis = append(h.chiis, hais)
	return nil
}

func (h *huroImpl) SetMinKan(hais [4]*hai.Hai) error {
	h.minKans = append(h.minKans, hais)
	return nil
}
func (h *huroImpl) SetAnKan(hais [4]*hai.Hai) error {
	h.ankans = append(h.ankans, hais)
	return nil
}

func (h *huroImpl) Kakan(inHai *hai.Hai) error {
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

	return HuroNotFoundErr
}
