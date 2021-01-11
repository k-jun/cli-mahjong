package tehai

import "mahjong/hai"

var _ Tehai = &TehaiMock{}

type TehaiMock struct {
	HaiMock   *hai.Hai
	HaisMock  []*hai.Hai
	ChiMock   [][2]*hai.Hai
	PonMock   [][2]*hai.Hai
	KanMock   [][3]*hai.Hai
	ErrorMock error
}

func (t *TehaiMock) Len() int {
	return len(t.HaisMock)
}

func (t *TehaiMock) Add(_ *hai.Hai) error {
	return t.ErrorMock
}

func (t *TehaiMock) Adds(_ []*hai.Hai) error {
	return t.ErrorMock
}

func (t *TehaiMock) Remove(_ *hai.Hai) (*hai.Hai, error) {
	return t.HaiMock, t.ErrorMock
}

func (t *TehaiMock) Removes(_ []*hai.Hai) ([]*hai.Hai, error) {
	return t.HaisMock, t.ErrorMock
}

func (t *TehaiMock) Replace(inHai *hai.Hai, _ *hai.Hai) (*hai.Hai, error) {
	outhai := t.HaiMock
	t.HaiMock = inHai
	return outhai, t.ErrorMock
}

func (t *TehaiMock) FindChiPairs(_ *hai.Hai) [][2]*hai.Hai {
	return t.ChiMock
}

func (t *TehaiMock) FindPonPairs(_ *hai.Hai) [][2]*hai.Hai {
	return t.PonMock
}

func (t *TehaiMock) FindKanPairs(_ *hai.Hai) [][3]*hai.Hai {
	return t.KanMock
}
