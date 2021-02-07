package yama

import "mahjong/model/hai"

var _ Yama = &YamaMock{}

type YamaMock struct {
	HaiMock   *hai.Hai
	ErrorMock error
	HaisMock  []*hai.Hai
}

func (y *YamaMock) Tumo() (*hai.Hai, error) {
	return y.HaiMock, y.ErrorMock
}

func (y *YamaMock) OmoteDora() []*hai.Hai {
	return y.HaisMock
}

func (y *YamaMock) UraDora() []*hai.Hai {
	return y.HaisMock
}

func (y *YamaMock) KanDora() error {
	return y.ErrorMock

}
