package tehai

import "mahjong/model/hai"

var _ Tehai = &TehaiMock{}

type TehaiMock struct {
	HaiMock    *hai.Hai
	HaisMock   []*hai.Hai
	ChiiMock   [][2]*hai.Hai
	PonMock    [][2]*hai.Hai
	MinKanMock [][3]*hai.Hai
	AnKanMock  [][4]*hai.Hai
	BoolMock   bool
	ErrorMock  error
}

func (t *TehaiMock) Hais() []*hai.Hai {
	return t.HaisMock
}

func (t *TehaiMock) Sort() error {
	return t.ErrorMock
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

func (t *TehaiMock) ChiiPairs(_ *hai.Hai) ([][2]*hai.Hai, error) {
	return t.ChiiMock, t.ErrorMock
}

func (t *TehaiMock) PonPairs(_ *hai.Hai) ([][2]*hai.Hai, error) {
	return t.PonMock, t.ErrorMock
}

func (t *TehaiMock) MinKanPairs(_ *hai.Hai) ([][3]*hai.Hai, error) {
	return t.MinKanMock, t.ErrorMock
}

func (t *TehaiMock) AnKanPairs(_ *hai.Hai) ([][4]*hai.Hai, error) {
	return t.AnKanMock, t.ErrorMock
}

func (t *TehaiMock) RiichiHais(_ *hai.Hai) ([]*hai.Hai, error) {
	return t.HaisMock, t.ErrorMock
}

func (t *TehaiMock) CanChii(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}

func (t *TehaiMock) CanPon(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}

func (t *TehaiMock) CanMinKan(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}

func (t *TehaiMock) CanAnKan(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}

func (t *TehaiMock) CanRiichi(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}

func (t *TehaiMock) CanRon(_ *hai.Hai) (bool, error) {
	return t.BoolMock, t.ErrorMock
}
