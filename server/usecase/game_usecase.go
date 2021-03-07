package usecase

import (
	"bytes"
	"fmt"
	"log"
	"mahjong/model/board"
	"mahjong/model/hai"
	"mahjong/model/player"
	"mahjong/model/view"
	"mahjong/storage"
	"regexp"
	"strconv"
	"strings"
)

type GameUsecase interface {
	JoinBoard(string, player.Player) (chan board.Board, error)
	InputController(string, player.Player)
	OutputController(string, player.Player, chan board.Board) error
}

type gameUsecaseImpl struct {
	BoardStorage storage.BoardStorage
	read         func([]byte) error
	write        func(string) error
}

var (
	re = regexp.MustCompile(`\d`)
)

func NewGameUsecase(ts storage.BoardStorage, write func(string) error, read func([]byte) error) GameUsecase {
	return &gameUsecaseImpl{
		BoardStorage: ts,
		read:         read,
		write:        write,
	}
}

func (gu *gameUsecaseImpl) InputController(id string, c player.Player) {
	Board, err := gu.BoardStorage.Find(id)
	if err != nil {
		log.Println(err)
	}

	for {
		buffer := make([]byte, 1024)
		if err := gu.read(buffer); err != nil {
			// dead check
			log.Println(err)
			if err := Board.LeavePlayer(c); err != nil {
				log.Println(err)
			}
			break
		}
		if string(buffer) != "" {
			buffer = bytes.Trim(buffer, "\x00")
			buffer = bytes.Trim(buffer, "\x10")
			haiName := strings.TrimSpace(string(buffer))
			turnIdx, err := Board.MyTurn(c)
			if err != nil {
				log.Println(err)
				Board.LeavePlayer(c)
				break
			}
			if Board.CurrentTurn() == turnIdx {
				// tsumo
				if haiName == "tsumo" {
					ok, err := c.CanTsumoAgari()
					if err != nil {
						log.Println(err)
						continue
					}
					if ok {
						if err := Board.SetWinIndex(turnIdx); err != nil {
							log.Println(err)
							continue
						}
						Board.Broadcast()
						fmt.Println(Board.LeavePlayer(c))
					}
				}
				// riichi or not
				if strings.HasPrefix(haiName, "riichi") {
					hais, err := c.Tehai().RiichiHais(c.Tsumohai())
					if err != nil {
						log.Println(err)
						continue
					}
					if len(hais) == 0 {
						log.Println(GameUsecaseInvalidActionErr)
						continue
					}
					idxstr := string(re.FindAll([]byte(haiName), 1)[0])
					idx, err := strconv.Atoi(idxstr)
					if idx >= len(hais) || idx < 0 {
						log.Println(GameUsecaseInvalidActionErr)
						continue
					}
					err = c.Riichi(hais[idx])
					if err != nil {
						log.Println(err)
						continue
					}
				} else {
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
				}
				err = Board.TurnEnd()
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				if strings.HasPrefix(haiName, "chii") {
					if Board.NextTurn() == turnIdx {
						err = Board.TakeAction(c, func(inHai *hai.Hai) error {
							pairs, err := c.Tehai().ChiiPairs(inHai)
							if err != nil {
								return err
							}
							if err != nil {
								return err
							}
							if len(pairs) == 0 {
								return GameUsecaseInvalidActionErr
							}
							idxstr := string(re.FindAll([]byte(haiName), 1)[0])
							idx, err := strconv.Atoi(idxstr)
							if idx >= len(pairs) || idx < 0 {
								return GameUsecaseInvalidActionErr
							}
							return c.Chii(inHai, pairs[idx])
						})
					}
				}
				if strings.HasPrefix(haiName, "pon") {
					err = Board.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().PonPairs(inHai)
						if err != nil {
							return err
						}
						if len(pairs) == 0 {
							return GameUsecaseInvalidActionErr
						}
						idxstr := string(re.FindAll([]byte(haiName), 1)[0])
						idx, err := strconv.Atoi(idxstr)
						if idx >= len(pairs) || idx < 0 {
							return GameUsecaseInvalidActionErr
						}
						return c.Pon(inHai, pairs[idx])
					})
				}
				if strings.HasPrefix(haiName, "kan") {
					err = Board.TakeAction(c, func(inHai *hai.Hai) error {
						pairs, err := c.Tehai().KanPairs(inHai)
						if err != nil {
							return err
						}
						if len(pairs) == 0 {
							return GameUsecaseInvalidActionErr
						}
						idxstr := string(re.FindAll([]byte(haiName), 1)[0])
						idx, err := strconv.Atoi(idxstr)
						if idx >= len(pairs) || idx < 0 {
							return GameUsecaseInvalidActionErr
						}
						return c.Kan(inHai, pairs[idx])
					})
				}
				if haiName == "ron" {
					err = Board.TakeAction(c, func(inHai *hai.Hai) error {
						isRon, err := c.Tehai().CanRon(inHai)
						if err != nil {
							return err
						}
						if isRon {
							if err := Board.SetWinIndex(turnIdx); err != nil {
								return err
							}
							Board.Broadcast()
						}
						return nil
					})
					// TODO where should i put this
					Board.LeavePlayer(c)
				}
				if haiName == "no" {
					err = Board.CancelAction(c)
				}
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

func (gu *gameUsecaseImpl) OutputController(id string, c player.Player, channel chan board.Board) error {
	_, err := gu.BoardStorage.Find(id)
	if err != nil {
		return err
	}
	for {
		Board, ok := <-channel
		if !ok {
			return GameUsecaseBoardChannelClosedErr
		}

		turnIdx, err := Board.MyTurn(c)
		if err != nil {
			return err
		}

		tehaistr, err := view.BoardString(c, Board)
		if err != nil {
			return err
		}

		// huros
		if Board.ActionCounter() != 0 && Board.CurrentTurn() != turnIdx {
			hai, err := Board.LastKawa()
			if err != nil {
				return err
			}
			actions := []player.Action{}
			// chii
			if Board.NextTurn() == turnIdx {
				chiis, err := c.Tehai().ChiiPairs(hai)
				if err != nil {
					return err
				}
				if len(chiis) != 0 {
					tehaistr += "\n" + "do you do Chii "
					for i, h := range chiis {
						tehaistr += strconv.Itoa(i) + ": " + h[0].Name() + "," + h[1].Name() + "   "
					}
				}
			}

			// pon
			pons, err := c.Tehai().PonPairs(hai)
			if err != nil {
				return err
			}
			if len(pons) != 0 {
				tehaistr += "\n" + "do you do Pon "
				for i, h := range pons {
					tehaistr += strconv.Itoa(i) + ": " + h[0].Name() + "," + h[1].Name() + "   "
				}
			}
			// kan
			kans, err := c.Tehai().KanPairs(hai)
			if err != nil {
				return err
			}
			if len(kans) != 0 {
				tehaistr += "\n" + "do you do Kan "
				for i, h := range kans {
					tehaistr += strconv.Itoa(i) + ": " + h[0].Name() + "," + h[1].Name() + "," + h[2].Name() + "   "
				}
			}
			// ron
			ok, err := c.Tehai().CanRon(hai)
			if err != nil {
				return err
			}
			if ok {
				actions = append(actions, player.Ron)
			}

			if err != nil {
				return err
			}
			for _, a := range actions {
				tehaistr += "\n" + "do you want " + string(a)
			}
		}
		// riichi
		hais, err := c.Tehai().RiichiHais(c.Tsumohai())
		if err != nil {
			return err
		}
		if len(hais) != 0 {
			tehaistr += "\n" + "do you do Riichi "
			for i, h := range hais {
				tehaistr += strconv.Itoa(i) + ": " + h.Name() + " "
			}
		}
		// tsumo agari
		ok, err = c.CanTsumoAgari()
		if err != nil {
			return err
		}
		if ok {
			tehaistr += "\n" + "do you do Tsumo "
		}
		if Board.ActionCounter() == 0 && Board.CurrentTurn() == turnIdx {
			tehaistr += "\n" + ">>"
		} else {
			tehaistr += "\n"
		}

		if err := gu.write(tehaistr); err != nil {
			return err
		}
	}
}

func (gu *gameUsecaseImpl) JoinBoard(id string, c player.Player) (chan board.Board, error) {
	Board, err := gu.BoardStorage.Find(id)
	if err != nil {
		return nil, err
	}

	return Board.JoinPlayer(c)
}
