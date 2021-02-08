package cha

import (
	"mahjong/model/hai"
	"mahjong/model/tehai"
	"mahjong/model/yama"
)

var _ Cha = &ChaMock{}

type ChaMock struct {
	ErrorMock error
	TehaiMock tehai.Tehai
	HaiMock   *hai.Hai
	HaisMock  []*hai.Hai
	BoolMock  bool
}

func (c *ChaMock) Tehai() tehai.Tehai {
	return c.TehaiMock
}

func (c *ChaMock) TumoHai() *hai.Hai {
	return c.HaiMock
}

func (c *ChaMock) Tumo() error {
	return c.ErrorMock
}

func (c *ChaMock) Dahai(outHai *hai.Hai) error {
	return c.ErrorMock
}
func (c *ChaMock) SetYama(_ yama.Yama) error {
	return c.ErrorMock
}

func (c *ChaMock) Haihai() error {
	return c.ErrorMock
}

func (c *ChaMock) Chi(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.ErrorMock
}

func (c *ChaMock) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.ErrorMock
}

func (c *ChaMock) Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	return c.ErrorMock
}

func (c *ChaMock) Kakan(inHai *hai.Hai) error {
	return c.ErrorMock
}

func (c *ChaMock) CanRichi() []*hai.Hai {
	return c.HaisMock
}

func (c *ChaMock) CanTumo() bool {
	return c.BoolMock

}
