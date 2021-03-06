package kawa

import (
	"mahjong/model/hai"
)

var _ Kawa = &KawaMock{}

type KawaMock struct {
	ErrorMock error
	HaiMock   *hai.Hai
	HaisMock  []*hai.Hai
}

func (h *KawaMock) Hais() []*hai.Hai {
	return h.HaisMock
}

func (h *KawaMock) Add(inHai *hai.Hai) error {
	h.HaiMock = inHai
	return h.ErrorMock
}

func (h *KawaMock) Last() (*hai.Hai, error) {
	return h.HaiMock, h.ErrorMock
}

func (h *KawaMock) RemoveLast() (*hai.Hai, error) {
	return h.HaiMock, h.ErrorMock
}
