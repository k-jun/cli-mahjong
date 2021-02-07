package huro

import (
	"mahjong/model/hai"
)

type Huro interface {
	Pon([3]*hai.Hai) error
	Chi([3]*hai.Hai) error
	MinKan([4]*hai.Hai) error
	AnKan([4]*hai.Hai) error
	Kakan(*hai.Hai) error
	GetPon() [][3]*hai.Hai
	GetChi() [][3]*hai.Hai
	GetMinKan() [][4]*hai.Hai
	GetAnKan() [][4]*hai.Hai
}

type huroImpl struct {
	pons    [][3]*hai.Hai
	chis    [][3]*hai.Hai
	minkans [][4]*hai.Hai
	ankans  [][4]*hai.Hai
}

func New() Huro {
	return &huroImpl{
		pons:    [][3]*hai.Hai{},
		chis:    [][3]*hai.Hai{},
		minkans: [][4]*hai.Hai{},
		ankans:  [][4]*hai.Hai{},
	}
}

func (h *huroImpl) GetPon() [][3]*hai.Hai {
	return h.pons
}
func (h *huroImpl) GetChi() [][3]*hai.Hai {
	return h.chis
}
func (h *huroImpl) GetMinKan() [][4]*hai.Hai {
	return h.minkans
}
func (h *huroImpl) GetAnKan() [][4]*hai.Hai {
	return h.ankans
}

func (h *huroImpl) Pon(hais [3]*hai.Hai) error {
	h.pons = append(h.pons, hais)
	return nil
}

func (h *huroImpl) Chi(hais [3]*hai.Hai) error {
	h.chis = append(h.chis, hais)
	return nil
}

func (h *huroImpl) MinKan(hais [4]*hai.Hai) error {
	h.minkans = append(h.minkans, hais)
	return nil
}
func (h *huroImpl) AnKan(hais [4]*hai.Hai) error {
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
			h.minkans = append(h.minkans, set)
			return nil
		}
	}

	return HuroNotFoundErr
}
