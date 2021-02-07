package huro

import "mahjong/model/hai"

var _ Huro = &HuroMock{}

type HuroMock struct {
	ErrorMock   error
	PonMock     [3]*hai.Hai
	ChiMock     [3]*hai.Hai
	MinKanMock  [4]*hai.Hai
	AnKanMock   [4]*hai.Hai
	PonsMock    [][3]*hai.Hai
	ChisMock    [][3]*hai.Hai
	MinKansMock [][4]*hai.Hai
	AnKansMock  [][4]*hai.Hai
}

func (h *HuroMock) GetPon() [][3]*hai.Hai {
	return h.PonsMock
}
func (h *HuroMock) GetChi() [][3]*hai.Hai {
	return h.ChisMock
}
func (h *HuroMock) GetMinKan() [][4]*hai.Hai {
	return h.MinKansMock
}
func (h *HuroMock) GetAnKan() [][4]*hai.Hai {
	return h.AnKansMock
}

func (h *HuroMock) Pon(hais [3]*hai.Hai) error {
	h.PonMock = hais

	return h.ErrorMock
}

func (h *HuroMock) Chi(hais [3]*hai.Hai) error {
	h.ChiMock = hais
	return h.ErrorMock
}

func (h *HuroMock) MinKan(hais [4]*hai.Hai) error {
	h.MinKanMock = hais
	return h.ErrorMock
}
func (h *HuroMock) AnKan(hais [4]*hai.Hai) error {
	h.AnKanMock = hais
	return h.ErrorMock
}

func (h *HuroMock) Kakan(x *hai.Hai) error {
	h.MinKanMock = [4]*hai.Hai{h.PonMock[0], h.PonMock[1], h.PonMock[2], x}
	h.PonMock = [3]*hai.Hai{}
	return h.ErrorMock
}
