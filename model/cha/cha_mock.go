package cha

import (
	"mahjong/model/hai"
	"mahjong/model/ho"
	"mahjong/model/huro"
	"mahjong/model/tehai"
	"mahjong/model/yama"
)

var _ Cha = &ChaMock{}

type ChaMock struct {
	ErrorMock       error
	TehaiMock       tehai.Tehai
	HuroMock        huro.Huro
	HaiMock         *hai.Hai
	HaisMock        []*hai.Hai
	HuroActionsMock []huro.HuroAction
	HoMock          ho.Ho
	BoolMock        bool
}

func (c *ChaMock) Tehai() tehai.Tehai {
	return c.TehaiMock
}

func (c *ChaMock) Ho() ho.Ho {
	return c.HoMock
}

func (c *ChaMock) Huro() huro.Huro {
	return c.HuroMock
}

func (c *ChaMock) Tsumohai() *hai.Hai {
	return c.HaiMock
}

func (c *ChaMock) Tsumo() error {
	return c.ErrorMock
}

func (c *ChaMock) Dahai(outHai *hai.Hai) error {
	return c.ErrorMock
}
func (c *ChaMock) SetYama(_ yama.Yama) error {
	return c.ErrorMock
}

func (c *ChaMock) Haipai() error {
	return c.ErrorMock
}

func (c *ChaMock) Chii(inHai *hai.Hai, outHais [2]*hai.Hai) error {
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

func (c *ChaMock) FindRiichiHai() ([]*hai.Hai, error) {
	return c.HaisMock, c.ErrorMock
}

func (c *ChaMock) CanTsumoAgari() (bool, error) {
	return c.BoolMock, c.ErrorMock

}
func (c *ChaMock) FindHuroActions(_ *hai.Hai) ([]huro.HuroAction, error) {
	return c.HuroActionsMock, c.ErrorMock
}
