package player

import (
	"mahjong/model/hai"
	"mahjong/model/kawa"
	"mahjong/model/naki"
	"mahjong/model/tehai"
	"mahjong/model/yama"
)

var _ Player = &PlayerMock{}

type PlayerMock struct {
	ErrorMock   error
	TehaiMock   tehai.Tehai
	NakiMock    naki.Naki
	HaiMock     *hai.Hai
	HaisMock    []*hai.Hai
	ActionsMock []Action
	KawaMock    kawa.Kawa
	BoolMock    bool
}

func (c *PlayerMock) Tehai() tehai.Tehai {
	return c.TehaiMock
}

func (c *PlayerMock) Kawa() kawa.Kawa {
	return c.KawaMock
}

func (c *PlayerMock) Naki() naki.Naki {
	return c.NakiMock
}

func (c *PlayerMock) Tsumohai() *hai.Hai {
	return c.HaiMock
}

func (c *PlayerMock) IsRiichi() bool {
	return c.BoolMock
}

func (c *PlayerMock) Tsumo() error {
	return c.ErrorMock
}

func (c *PlayerMock) Dahai(outHai *hai.Hai) error {
	return c.ErrorMock
}

func (c *PlayerMock) Riichi(_ *hai.Hai) error {
	return c.ErrorMock
}

func (c *PlayerMock) SetYama(_ yama.Yama) error {
	return c.ErrorMock
}

func (c *PlayerMock) Haipai() error {
	return c.ErrorMock
}

func (c *PlayerMock) Chii(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.ErrorMock
}

func (c *PlayerMock) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	return c.ErrorMock
}

func (c *PlayerMock) MinKan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	return c.ErrorMock
}

func (c *PlayerMock) AnKan(_ [4]*hai.Hai) error {
	return c.ErrorMock
}
func (c *PlayerMock) Kakan() error {
	return c.ErrorMock
}

func (c *PlayerMock) FindRiichiHai() ([]*hai.Hai, error) {
	return c.HaisMock, c.ErrorMock
}

func (c *PlayerMock) CanTsumoAgari() (bool, error) {
	return c.BoolMock, c.ErrorMock
}

func (c *PlayerMock) CanRiichi() (bool, error) {
	return c.BoolMock, c.ErrorMock
}

func (c *PlayerMock) CanAnKan() (bool, error) {
	return c.BoolMock, c.ErrorMock
}

func (c *PlayerMock) Actions(_ *hai.Hai) ([]Action, error) {
	return c.ActionsMock, c.ErrorMock
}
