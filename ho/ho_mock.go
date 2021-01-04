package ho

import (
	"mahjong/hai"
)

var _ Ho = &HoMock{}

type HoMock struct {
	ErrorMock error
	HaiMock   *hai.Hai
}

func (h *HoMock) Add(inHai *hai.Hai) error {
	h.HaiMock = inHai
	return h.ErrorMock
}

func (h *HoMock) Last() (*hai.Hai, error) {
	return h.HaiMock, h.ErrorMock
}
