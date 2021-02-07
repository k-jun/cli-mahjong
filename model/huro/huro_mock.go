package huro

import "mahjong/model/hai"

var _ Huro = &HuroMock{}

type HuroMock struct {
	ErrorMock error
	PonMock   [3]*hai.Hai
	ChiMock   [3]*hai.Hai
	KanMock   [4]*hai.Hai
}

func (h *HuroMock) Pon(hais [3]*hai.Hai) error {
	h.PonMock = hais

	return h.ErrorMock
}

func (h *HuroMock) Chi(hais [3]*hai.Hai) error {
	h.ChiMock = hais
	return h.ErrorMock
}

func (h *HuroMock) Kan(hais [4]*hai.Hai) error {
	h.KanMock = hais
	return h.ErrorMock
}

func (h *HuroMock) Kakan(x *hai.Hai) error {
	h.KanMock = [4]*hai.Hai{h.PonMock[0], h.PonMock[1], h.PonMock[2], x}
	h.PonMock = [3]*hai.Hai{}
	return h.ErrorMock
}
