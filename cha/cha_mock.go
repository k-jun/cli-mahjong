package cha

import (
	"mahjong/hai"
	"mahjong/tehai"
	"mahjong/yama"
)

var _ Cha = &ChaMock{}

type ChaMock struct {
	ErrorMock error
	TehaiMock tehai.Tehai
}

func (c *ChaMock) Tehai() tehai.Tehai {
	return c.TehaiMock
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
