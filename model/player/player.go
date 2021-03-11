package player

import (
	"mahjong/model/hai"
	"mahjong/model/kawa"
	"mahjong/model/naki"
	"mahjong/model/tehai"
	"mahjong/model/yama"

	"github.com/google/uuid"
)

type Player interface {
	// getter
	Tehai() tehai.Tehai
	Kawa() kawa.Kawa
	Tsumohai() *hai.Hai
	Naki() naki.Naki
	IsRiichi() bool
	// setter
	SetYama(yama.Yama) error

	// judger
	CanRiichi() (bool, error)
	CanTsumoAgari() (bool, error)
	CanAnKan() (bool, error)

	Tsumo() error
	Dahai(*hai.Hai) error
	Haipai() error
	Chii(*hai.Hai, [2]*hai.Hai) error
	Pon(*hai.Hai, [2]*hai.Hai) error
	AnKan([4]*hai.Hai) error
	MinKan(*hai.Hai, [3]*hai.Hai) error
	Kakan() error
	Riichi(*hai.Hai) error
}

type Action string

var (
	Chii       Action = "chii"
	Pon        Action = "pon"
	Kan        Action = "kan"
	Ron        Action = "ron"
	Riichi     Action = "riichi"
	Tsumo      Action = "tsumo"
	AllActions        = []Action{Chii, Pon, Kan, Ron, Riichi, Tsumo}
)

func AtoAction(s string) (Action, error) {
	for _, a := range AllActions {
		if string(a) == s {
			return a, nil
		}
	}

	return "", PlayerActionInvalidErr
}

type playerImpl struct {
	id       uuid.UUID
	tsumohai *hai.Hai
	kawa     kawa.Kawa
	tehai    tehai.Tehai
	naki     naki.Naki
	yama     yama.Yama
	isRiichi bool
}

func New(id uuid.UUID, k kawa.Kawa, t tehai.Tehai, n naki.Naki) Player {
	return &playerImpl{
		id:       id,
		kawa:     k,
		tehai:    t,
		naki:     n,
		yama:     nil,
		isRiichi: false,
	}
}

func (c *playerImpl) Tehai() tehai.Tehai {
	return c.tehai
}

func (c *playerImpl) Kawa() kawa.Kawa {
	return c.kawa
}

func (c *playerImpl) Naki() naki.Naki {
	return c.naki
}

func (c *playerImpl) Tsumohai() *hai.Hai {
	return c.tsumohai
}

func (c *playerImpl) IsRiichi() bool {
	return c.isRiichi
}

func (c *playerImpl) Tsumo() error {
	if c.tsumohai != nil {
		return PlayerAlreadyHaveTsumohaiErr
	}

	tsumohai, err := c.yama.Draw()
	if err != nil {
		return err
	}

	c.tsumohai = tsumohai
	return nil
}

func (c *playerImpl) Dahai(outHai *hai.Hai) error {
	var err error
	if c.isRiichi && outHai != c.tsumohai {
		return PlayerAlreadyRiichiErr
	}
	if outHai != c.tsumohai {
		if c.tsumohai == nil {
			outHai, err = c.tehai.Remove(outHai)
		} else {
			outHai, err = c.tehai.Replace(c.tsumohai, outHai)
		}

		if err != nil {
			return err
		}
		if err := c.tehai.Sort(); err != nil {
			return err
		}
	}
	c.tsumohai = nil

	return c.kawa.Add(outHai)
}

func (c *playerImpl) SetYama(y yama.Yama) error {
	if c.yama != nil {
		return PlayerAlreadyHaveYamaErr
	}
	c.yama = y
	return nil
}

func (c *playerImpl) Haipai() error {
	if len(c.tehai.Hais()) != 0 {
		return PlayerAlreadyDidHaipaiErr
	}

	for i := 0; i < tehai.MaxHaisLen; i++ {
		tsumoHai, err := c.yama.Draw()
		if err != nil {
			return err
		}
		if err := c.tehai.Add(tsumoHai); err != nil {
			return err
		}
	}

	if err := c.tehai.Sort(); err != nil {
		return err
	}
	return nil
}

func (c *playerImpl) Chii(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.naki.SetChii(meld)
}

func (c *playerImpl) Pon(inHai *hai.Hai, outHais [2]*hai.Hai) error {
	if inHai == c.tsumohai {
		c.tsumohai = nil
	}
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1]})
	if err != nil {
		return err
	}
	meld := [3]*hai.Hai{inHai, hais[0], hais[1]}

	return c.naki.SetPon(meld)
}

func (c *playerImpl) MinKan(inHai *hai.Hai, outHais [3]*hai.Hai) error {
	hais, err := c.tehai.Removes([]*hai.Hai{outHais[0], outHais[1], outHais[2]})
	if err != nil {
		return err
	}
	meld := [4]*hai.Hai{inHai, hais[0], hais[1], hais[2]}

	return c.naki.SetMinKan(meld)
}

func (c *playerImpl) AnKan(hais [4]*hai.Hai) error {
	if c.tsumohai == hais[0] {
		if _, err := c.tehai.Removes([]*hai.Hai{hais[0], hais[1], hais[2]}); err != nil {
			return err
		}
	} else {
		if _, err := c.tehai.Removes([]*hai.Hai{hais[0], hais[1], hais[2], hais[3]}); err != nil {
			return err
		}
		if err := c.tehai.Add(c.tsumohai); err != nil {
			return err
		}
	}
	c.tsumohai = nil
	return c.naki.SetAnKan([4]*hai.Hai{hais[0], hais[1], hais[2], hais[3]})
}

func (c *playerImpl) Kakan() error {
	err := c.naki.Kakan(c.tsumohai)
	if err != nil {
		return err
	}
	c.tsumohai = nil
	return nil
}

func (c *playerImpl) Riichi(inHai *hai.Hai) error {
	if c.isRiichi {
		return PlayerAlreadyRiichiErr
	}
	err := c.Dahai(inHai)
	if err != nil {
		return err
	}
	c.isRiichi = true
	return nil
}

func (c *playerImpl) CanAnKan() (bool, error) {
	return c.tehai.CanAnKan(c.tsumohai)
}

func (c *playerImpl) CanRiichi() (bool, error) {
	cntChii := len(c.naki.Chiis())
	cntPon := len(c.naki.Pons())
	cntKan := len(c.naki.MinKans())
	if cntChii == 0 && cntPon == 0 && cntKan == 0 {
		return c.tehai.CanRiichi(c.tsumohai)
	}

	return false, nil
}

func (c *playerImpl) CanTsumoAgari() (bool, error) {
	return c.tehai.CanRon(c.tsumohai)
}
