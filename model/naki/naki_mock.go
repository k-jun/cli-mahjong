package naki

import "mahjong/model/hai"

var _ Naki = &NakiMock{}

type NakiMock struct {
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

func (h *NakiMock) Pons() [][3]*hai.Hai {
	return h.PonsMock
}
func (h *NakiMock) Chiis() [][3]*hai.Hai {
	return h.ChiisMock
}
func (h *NakiMock) MinKans() [][4]*hai.Hai {
	return h.MinKansMock
}
func (h *NakiMock) AnKans() [][4]*hai.Hai {
	return h.AnKansMock
}

func (h *NakiMock) SetPon(hais [3]*hai.Hai) error {
	h.PonMock = hais
	return h.ErrorMock
}

func (h *NakiMock) SetChii(hais [3]*hai.Hai) error {
	h.ChiiMock = hais
	return h.ErrorMock
}

func (h *NakiMock) SetMinKan(hais [4]*hai.Hai) error {
	h.MinKanMock = hais
	return h.ErrorMock
}
func (h *NakiMock) SetAnKan(hais [4]*hai.Hai) error {
	h.AnKanMock = hais
	return h.ErrorMock
}

func (h *NakiMock) Kakan(x *hai.Hai) error {
	h.MinKanMock = [4]*hai.Hai{h.PonMock[0], h.PonMock[1], h.PonMock[2], x}
	h.PonMock = [3]*hai.Hai{}
	return h.ErrorMock
}
