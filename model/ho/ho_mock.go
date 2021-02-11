package ho

import (
	"mahjong/model/hai"
)

var _ Ho = &HoMock{}

type HoMock struct {
	ErrorMock error
	HaiMock   *hai.Hai
	HaisMock  []*hai.Hai
}

func (h *HoMock) Hais() []*hai.Hai {
	return h.HaisMock
}

func (h *HoMock) Add(inHai *hai.Hai) error {
	h.HaiMock = inHai
	return h.ErrorMock
}

func (h *HoMock) Last() (*hai.Hai, error) {
	return h.HaiMock, h.ErrorMock
}

func (h *HoMock) RemoveLast() (*hai.Hai, error) {
	return h.HaiMock, h.ErrorMock
}
