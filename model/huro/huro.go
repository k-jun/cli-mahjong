package huro

import (
	"mahjong/model/hai"
)

type Huro interface {
	Pon([3]*hai.Hai) error
	Chi([3]*hai.Hai) error
	Kan([4]*hai.Hai) error
	Kakan(*hai.Hai) error
}

type huroImpl struct {
	pons [][3]*hai.Hai
	chis [][3]*hai.Hai
	kans [][4]*hai.Hai
}

func New() Huro {
	return &huroImpl{
		pons: [][3]*hai.Hai{},
		chis: [][3]*hai.Hai{},
		kans: [][4]*hai.Hai{},
	}

}

func (h *huroImpl) Pon(hais [3]*hai.Hai) error {
	h.pons = append(h.pons, hais)
	return nil
}

func (h *huroImpl) Chi(hais [3]*hai.Hai) error {
	h.chis = append(h.chis, hais)
	return nil
}

func (h *huroImpl) Kan(hais [4]*hai.Hai) error {
	h.kans = append(h.kans, hais)
	return nil
}

func (h *huroImpl) Kakan(inHai *hai.Hai) error {
	for idx, pon := range h.pons {
		if pon[0] == inHai {
			h.pons[idx] = h.pons[0]
			h.pons = h.pons[1:]

			set := [4]*hai.Hai{}
			set[0], set[1], set[2], set[3] = pon[0], pon[1], pon[2], inHai
			h.kans = append(h.kans, set)
			return nil
		}
	}

	return HuroNotFoundErr
}
