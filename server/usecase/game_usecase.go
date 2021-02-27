package usecase

import (
	"bytes"
	"fmt"
	"log"
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/huro"
	"mahjong/model/taku"
	"mahjong/storage"
	"strings"
)

type GameUsecase interface {
	JoinTaku(string, cha.Cha) (chan taku.Taku, error)
	InputController(string, cha.Cha)
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

func (gu *gameUsecaseImpl) InputController(id string, c cha.Cha) {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		log.Println(err)
	}

	for {
		buffer := make([]byte, 1024)
		if err := gu.read(buffer); err != nil {
			// dead check
			log.Println(err)
			if err := taku.LeaveCha(c); err != nil {
				log.Println(err)
			}
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
				if haiName == "tsumo" {
					ok, err := c.CanTsumoAgari()
					if err != nil {
						log.Println(err)
						continue
					}
					if ok {
						if err := taku.SetWinIndex(turnIdx); err != nil {
							log.Println(err)
							continue
						}
						taku.Broadcast()
						fmt.Println(taku.LeaveCha(c))
					}
				}
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
					if taku.NextTurn() == turnIdx {
						err = taku.TakeAction(c, func(inHai *hai.Hai) error {
							pairs, err := c.Tehai().FindChiiPairs(inHai)
							if err != nil {
								return err
							}
							if len(pairs) == 0 {
								return GameUsecaseHuroNotFoundErr
							}
							return c.Chii(inHai, pairs[0])
						})
					}
				}
				if haiName == "pon" {
					err = taku.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().FindPonPairs(inHai)
						if err != nil {
							return err
						}
						if len(pairs) == 0 {
							return GameUsecaseHuroNotFoundErr
						}
						return c.Pon(inHai, pairs[0])
					})
				}
				if haiName == "kan" {
					err = taku.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().FindKanPairs(inHai)
						if err != nil {
							return err
						}
						if len(pairs) == 0 {
							return GameUsecaseHuroNotFoundErr
						}
						return c.Kan(inHai, pairs[0])
					})
				}
				if haiName == "ron" {
					err = taku.TakeAction(c, func(inHai *hai.Hai) error {
						isRon, err := c.CanRon(inHai)
						if err != nil {
							return err
						}
						if isRon {
							if err := taku.SetWinIndex(turnIdx); err != nil {
								return err
							}
							taku.Broadcast()
						}
						return nil
					})
					// TODO where should i put this
					taku.LeaveCha(c)
				}
				if haiName == "no" {
					err = taku.CancelAction(c)
				}
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

func (gu *gameUsecaseImpl) OutputController(id string, c cha.Cha, channel chan taku.Taku) error {
	_, err := gu.takuStorage.Find(id)
	if err != nil {
		return err
	}
	for {
		taku, ok := <-channel
		if !ok {
			return GameUsecaseTakuChannelClosedErr
		}

		turnIdx, err := taku.MyTurn(c)
		if err != nil {
			return err
		}

		tehaistr := taku.Draw(c)

		// huros
		if taku.ActionCounter() != 0 && taku.CurrentTurn() != turnIdx {
			hai, err := taku.LastHo()
			if err != nil {
				return err
			}
			actions := []huro.HuroAction{}
			// chii
			if taku.NextTurn() == turnIdx {
				chiis, err := c.Tehai().FindChiiPairs(hai)
				if err != nil {
					return err
				}
				if len(chiis) != 0 {
					actions = append(actions, huro.Chii)
				}
			}

			// pon
			pons, err := c.Tehai().FindPonPairs(hai)
			if err != nil {
				return err
			}
			if len(pons) != 0 {
				actions = append(actions, huro.Pon)
			}
			// kan
			kans, err := c.Tehai().FindKanPairs(hai)
			if err != nil {
				return err
			}
			if len(kans) != 0 {
				actions = append(actions, huro.Kan)
			}
			// ron
			ok, err := c.CanRon(hai)
			if err != nil {
				return err
			}
			if ok {
				actions = append(actions, huro.Ron)
			}

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
			tehaistr += "\n" + "do you do riichi?: "
		}
		// tsumo agari
		ok, err = c.CanTsumoAgari()
		if err != nil {
			return err
		}
		if ok {
			tehaistr += "\n" + "do you do tsumo?: "
		}
		if taku.ActionCounter() == 0 && taku.CurrentTurn() == turnIdx {
			tehaistr += "\n" + ">>"
		} else {
			tehaistr += "\n"
		}

		if err := gu.write(tehaistr); err != nil {
			return err
		}
	}
}

func (gu *gameUsecaseImpl) JoinTaku(id string, c cha.Cha) (chan taku.Taku, error) {
	taku, err := gu.takuStorage.Find(id)
	if err != nil {
		return nil, err
	}

	return taku.JoinCha(c)
}
