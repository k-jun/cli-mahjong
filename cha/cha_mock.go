package cha

import "mahjong/hai"

var _ Cha = &ChaMock{}

type ChaMock struct {
	OutError error
}

func (c *ChaMock) Tumo() error {
	return c.OutError
}

func (c *ChaMock) Dahai(outHai *hai.Hai) error {
	return c.OutError
}

func (c *ChaMock) Chi(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.OutError
}

func (c *ChaMock) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.OutError
}

func (c *ChaMock) Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	return c.OutError
}

func (c *ChaMock) Kakan(inHai *hai.Hai) error {
	return c.OutError
}
