package cha

import (
	"mahjong/hai"
	"mahjong/ho"
	"mahjong/huro"
	"mahjong/tehai"
	"mahjong/yama"

	"github.com/google/uuid"
)

type Cha interface {
	Tumo() error
	Dahai(outHai *hai.Hai) error
	Chi(inHai *hai.Hai, outHais [2]*hai.Hai) error
	Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error
	Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error
	Kakan(inHai *hai.Hai) error
}

type chaImpl struct {
	id      uuid.UUID
	tumohai *hai.Hai
	ho      ho.Ho
	tehai   tehai.Tehai
	huro    huro.Huro
	yama    yama.Yama
}

func New(id uuid.UUID, ho ho.Ho, t tehai.Tehai, y yama.Yama, hu huro.Huro) Cha {
	return &chaImpl{
		id:    id,
		ho:    ho,
		tehai: t,
		yama:  y,
		huro:  hu,
	}
}

func (c *chaImpl) Id() uuid.UUID {
	return c.id
}

func (c *chaImpl) Tumo() error {
	if c.tumohai != nil {
		return ChaAlreadyHaveTumohaiErr
	}

	tumohai, err := c.yama.Tumo()
	if err != nil {
		return err
	}

	c.tumohai = tumohai
	return nil
}

func (c *chaImpl) Dahai(outHai *hai.Hai) error {
	var err error
	if outHai != c.tumohai {
		outHai, err = c.tehai.Replace(c.tumohai, outHai)
		if err != nil {
			return err
		}
	}
	c.tumohai = nil

	return c.ho.Add(outHai)
}

func (c *chaImpl) Chi(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tumohai {
		c.tumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	set := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.huro.Chi(set)
}

func (c *chaImpl) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tumohai {
		c.tumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	set := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.huro.Pon(set)
}

func (c *chaImpl) Kan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	if inHai == c.tumohai {
		c.tumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	set := [4]*hai.Hai{inHai, hais[0], hais[1], hais[2]}

	return c.huro.Kan(set)
}

func (c *chaImpl) Kakan(inHai *hai.Hai) error {
	if inHai == c.tumohai {
		c.tumohai = nil
	}
	return c.huro.Kakan(inHai)
}
