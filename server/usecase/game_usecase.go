package usecase

import (
	"bytes"
	"fmt"
	"log"
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/taku"
	"mahjong/storage"
	"strings"
)

type GameUsecase interface {
	JoinTaku(string, cha.Cha) (chan taku.Taku, error)
	InputController(string, cha.Cha) error
	OutputController(string, cha.Cha, chan taku.Taku) error
}

type gameUsecaseImpl struct {
	takuStorage storage.TakuStorage
	read        func([]byte) error
	write       func(string) error
}

func NewGameUsecase(ts storage.TakuStorage, write func(string) error, read func([]byte) error) GameUsecase {
	return &gameUsecaseImpl{
		takuStorage: ts,
		read:        read,
		write:       write,
	}
}

func (gu *gameUsecaseImpl) InputController(id string, c cha.Cha) error {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return err
	}

	for {
		buffer := make([]byte, 1024)
		if err := gu.read(buffer); err != nil {
			// dead check
			taku.LeaveCha(c)
			break
		}
		if string(buffer) != "" {
			buffer = bytes.Trim(buffer, "\x00")
			buffer = bytes.Trim(buffer, "\x10")
			haiName := strings.TrimSpace(string(buffer))
			turnIdx, err := taku.MyTurn(c)
			if err != nil {
				log.Println(err)
				taku.LeaveCha(c)
				break
			}
			if taku.CurrentTurn() == turnIdx {
				outHai, err := hai.AtoHai(haiName)
				if err != nil {
					log.Println(err)
					continue
				}
				err = c.Dahai(outHai)
				if err != nil {
					log.Println(err)
					continue
				}
				err = taku.TurnEnd()
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				if haiName == "chii" {
					taku.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().FindChiiPairs(inHai)
						if err != nil {
							return err
						}
						return c.Chii(inHai, pairs[0])
					})
				}
				if haiName == "pon" {
					taku.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().FindPonPairs(inHai)
						if err != nil {
							return err
						}
						return c.Pon(inHai, pairs[0])
					})
				}
				if haiName == "kan" {
					taku.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().FindKanPairs(inHai)
						if err != nil {
							return err
						}
						return c.Kan(inHai, pairs[0])
					})
				}
				if haiName == "ron" {
					taku.TakeAction(c, func(inHai *hai.Hai) error {
						isRon, err := c.CanRon(inHai)
						if err != nil {
							return err
						}
						if isRon {
							fmt.Println("TODO: Ron! GAME END")
						}
						return nil
					})
				}
				if haiName == "no" {
					taku.CancelAction(c)
				}
			}
		}
	}

	return nil
}

func (gu *gameUsecaseImpl) OutputController(id string, c cha.Cha, channel chan taku.Taku) error {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return err
	}
	for {
		_, ok := <-channel
		if ok {
			break
		}

		turnIdx, err := taku.MyTurn(c)
		if err != nil {
			log.Println(err)
			return err
		}

		tehaistr := taku.Draw(c)

		// huros
		if taku.ActionCounter() != 0 && taku.CurrentTurn() != turnIdx {
			hai, err := taku.LastHo()
			if err != nil {
				return err
			}
			actions, err := c.FindHuroActions(hai)
			if err != nil {
				return err
			}
			for _, a := range actions {
				tehaistr += "\n" + "do you want " + string(a)
			}
		}
		// riichi
		hais, err := c.FindRiichiHai()
		if err != nil {
			return err
		}
		if len(hais) != 0 {
			tehaistr += "\n" + "do you do richi?: "
		}
		// tsumo agari
		ok, err = c.CanTsumoAgari()
		if err != nil {
			return err
		}
		if ok {
			tehaistr += "\n" + "do you do tumo?: "
		}
		tehaistr += `
		`
		gu.write(tehaistr + "\n")
	}

	return nil
}

func (gu *gameUsecaseImpl) JoinTaku(id string, c cha.Cha) (chan taku.Taku, error) {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return nil, err
	}

	return taku.JoinCha(c)
}
