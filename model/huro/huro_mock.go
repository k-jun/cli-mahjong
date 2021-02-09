package huro

import "mahjong/model/hai"

var _ Huro = &HuroMock{}

type HuroMock struct {
	ErrorMock error

	PonMock     [3]*hai.Hai
	PonsMock    [][3]*hai.Hai
	ChiiMock    [3]*hai.Hai
	ChiisMock   [][3]*hai.Hai
	MinKanMock  [4]*hai.Hai
	MinKansMock [][4]*hai.Hai
	AnKanMock   [4]*hai.Hai
	AnKansMock  [][4]*hai.Hai
}

func (h *HuroMock) Pons() [][3]*hai.Hai {
	return h.PonsMock
}
func (h *HuroMock) Chiis() [][3]*hai.Hai {
	return h.ChiisMock
}
func (h *HuroMock) MinKans() [][4]*hai.Hai {
	return h.MinKansMock
}
func (h *HuroMock) AnKans() [][4]*hai.Hai {
	return h.AnKansMock
}

func (h *HuroMock) SetPon(hais [3]*hai.Hai) error {
	h.PonMock = hais
	return h.ErrorMock
}

func (h *HuroMock) SetChii(hais [3]*hai.Hai) error {
	h.ChiiMock = hais
	return h.ErrorMock
}

func (h *HuroMock) SetMinKan(hais [4]*hai.Hai) error {
	h.MinKanMock = hais
	return h.ErrorMock
}
func (h *HuroMock) SetAnKan(hais [4]*hai.Hai) error {
	h.AnKanMock = hais
	return h.ErrorMock
}

func (h *HuroMock) Kakan(x *hai.Hai) error {
	h.MinKanMock = [4]*hai.Hai{h.PonMock[0], h.PonMock[1], h.PonMock[2], x}
	h.PonMock = [3]*hai.Hai{}
	return h.ErrorMock
}
