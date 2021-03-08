package usecase

import (
	"bytes"
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
	OutputController(player.Player, chan board.Board) error
}

type gameUsecaseImpl struct {
	BoardStorage storage.BoardStorage
	read         func([]byte) error
	write        func(string) error
}

var (
	num = regexp.MustCompile(`\d+`)
	str = regexp.MustCompile(`\w+`)
)

func NewGameUsecase(ts storage.BoardStorage, write func(string) error, read func([]byte) error) GameUsecase {
	return &gameUsecaseImpl{
		BoardStorage: ts,
		read:         read,
		write:        write,
	}
}

type ActionType string

var (
	Normal ActionType = ""
	Tsumo  ActionType = "tsumo"
	Riichi ActionType = "riichi"
	Chii   ActionType = "chii"
	Pon    ActionType = "pon"
	Kan    ActionType = "kan"
	Ron    ActionType = "ron"
	Cancel ActionType = "no"
)

type InputCommand struct {
	actionType  ActionType
	actionIndex int
	hai         *hai.Hai
}

func (gu *gameUsecaseImpl) CommandParser(raw []byte) (*InputCommand, error) {
	raw = bytes.Trim(raw, "\x00")
	raw = bytes.Trim(raw, "\x10")
	rawstr := strings.TrimSpace(string(raw))
	h, err := hai.AtoHai(rawstr)
	if err != nil && err != hai.HaiInvalidArgumentErr {
		return nil, err
	}

	ic := InputCommand{actionType: Normal, actionIndex: 0, hai: h}
	if h != nil {
		return &ic, nil
	}

	actions := str.FindAllString(rawstr, 1)
	if len(actions) != 1 {
		return &ic, GameUsecaseInvalidActionErr
	}

	switch actions[0] {
	case "tsumo":
		ic.actionType = Tsumo
	case "riichi":
		ic.actionType = Riichi
	case "chii":
		ic.actionType = Chii
	case "pon":
		ic.actionType = Pon
	case "kan":
		ic.actionType = Kan
	case "ron":
		ic.actionType = Ron
	case "no":
		ic.actionType = Cancel
	default:
		return &ic, GameUsecaseInvalidActionErr
	}
	idxes := num.FindAllString(rawstr, 1)
	if len(idxes) != 1 {
		return &ic, nil
	}
	idx, _ := strconv.Atoi(idxes[0])
	ic.actionIndex = idx

	return &ic, nil
}

func (gu *gameUsecaseImpl) Normal(b board.Board, p player.Player, ic *InputCommand) error {
	err := p.Dahai(ic.hai)
	if err != nil {
		return err
	}
	return b.TurnEnd()
}

func (gu *gameUsecaseImpl) Tsumo(b board.Board, p player.Player, ic *InputCommand) error {
	ok, err := p.CanTsumoAgari()
	if err != nil {
		return err
	}

	if !ok {
		return GameUsecaseInvalidActionErr
	}
	idx, err := b.MyTurn(p)
	if err != nil {
		return err
	}
	if err := b.SetWinIndex(idx); err != nil {
		return err
	}
	b.Broadcast()
	return b.LeavePlayer(p)
}

func (gu *gameUsecaseImpl) Riichi(b board.Board, p player.Player, ic *InputCommand) error {
	hais, err := p.Tehai().RiichiHais(p.Tsumohai())
	if err != nil {
		return err
	}
	if len(hais) == 0 {
		return GameUsecaseInvalidActionErr
	}

	if ic.actionIndex >= len(hais) || ic.actionIndex < 0 {
		return GameUsecaseInvalidActionErr
	}
	err = p.Riichi(hais[ic.actionIndex])
	if err != nil {
		log.Println(err)
	}
	return b.TurnEnd()
}

func (gu *gameUsecaseImpl) AnKan(b board.Board, p player.Player, ic InputCommand) {

}

func (gu *gameUsecaseImpl) Chii(b board.Board, p player.Player, ic *InputCommand) error {
	return b.TakeAction(p, func(inHai *hai.Hai) error {
		pairs, err := p.Tehai().ChiiPairs(inHai)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if len(pairs) == 0 {
			return GameUsecaseInvalidActionErr
		}
		if ic.actionIndex >= len(pairs) || ic.actionIndex < 0 {
			return GameUsecaseInvalidActionErr
		}
		return p.Chii(inHai, pairs[ic.actionIndex])
	})
}

func (gu *gameUsecaseImpl) Pon(b board.Board, p player.Player, ic *InputCommand) error {
	return b.TakeAction(p, func(inHai *hai.Hai) error {
		pairs, err := p.Tehai().PonPairs(inHai)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if len(pairs) == 0 {
			return GameUsecaseInvalidActionErr
		}
		if ic.actionIndex >= len(pairs) || ic.actionIndex < 0 {
			return GameUsecaseInvalidActionErr
		}
		return p.Pon(inHai, pairs[ic.actionIndex])
	})
}

func (gu *gameUsecaseImpl) MinKan(b board.Board, p player.Player, ic *InputCommand) error {
	return b.TakeAction(p, func(inHai *hai.Hai) error {
		pairs, err := p.Tehai().MinKanPairs(inHai)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if len(pairs) == 0 {
			return GameUsecaseInvalidActionErr
		}
		if ic.actionIndex >= len(pairs) || ic.actionIndex < 0 {
			return GameUsecaseInvalidActionErr
		}
		return p.MinKan(inHai, pairs[ic.actionIndex])
	})
}

func (gu *gameUsecaseImpl) Ron(b board.Board, p player.Player, ic *InputCommand) error {
	turnIdx, err := b.MyTurn(p)
	if err != nil {
		return err
	}

	err = b.TakeAction(p, func(inHai *hai.Hai) error {
		ok, err := p.Tehai().CanRon(inHai)
		if err != nil {
			return err
		}
		if !ok {
			return GameUsecaseInvalidActionErr
		}
		if err := b.SetWinIndex(turnIdx); err != nil {
			return err
		}
		b.Broadcast()
		return nil
	})
	if err != nil {
		return err
	}
	return b.LeavePlayer(p)
}

func (gu *gameUsecaseImpl) InputController(id string, p player.Player) {
	b, err := gu.BoardStorage.Find(id)
	if err != nil {
		log.Println(err)
	}

	for {
		buffer := make([]byte, 1024)
		if err := gu.read(buffer); err != nil {
			// dead check
			log.Println(err)
			if err := b.LeavePlayer(p); err != nil {
				log.Println(err)
			}
			break
		}

		command, err := gu.CommandParser(buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		turnIdx, err := b.MyTurn(p)
		if err != nil {
			log.Println(err)
			b.LeavePlayer(p)
			break
		}
		if b.CurrentTurn() == turnIdx {
			// my turn
			switch command.actionType {
			case Normal:
				err = gu.Normal(b, p, command)
			case Tsumo:
				err = gu.Tsumo(b, p, command)
			case Riichi:
				err = gu.Riichi(b, p, command)
			case Kan:
				continue
			default:
				continue
			}
		} else {
			// not my turn
			if b.NextTurn() == turnIdx && command.actionType == Chii {
				err = gu.Chii(b, p, command)
			} else {
				switch command.actionType {
				case Pon:
					err = gu.Pon(b, p, command)
				case Kan:
					err = gu.MinKan(b, p, command)
				case Ron:
					err = gu.Ron(b, p, command)
				case Cancel:
					err = b.CancelAction(p)
				default:
					continue
				}
			}
		}
		if err != nil {
			log.Println(err)
		}
	}
}

func (gu *gameUsecaseImpl) OutputController(c player.Player, channel chan board.Board) error {
	for {
		board, ok := <-channel
		if !ok {
			return GameUsecaseBoardChannelClosedErr
		}

		turnIdx, err := board.MyTurn(c)
		if err != nil {
			return err
		}

		tehaistr, err := view.BoardString(c, board)
		if err != nil {
			return err
		}

		// huros
		if board.ActionCounter() != 0 && board.CurrentTurn() != turnIdx {
			hai, err := board.LastKawa()
			if err != nil {
				return err
			}
			actions := []player.Action{}
			// chii
			if board.NextTurn() == turnIdx {
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
			kans, err := c.Tehai().MinKanPairs(hai)
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
		if board.ActionCounter() == 0 && board.CurrentTurn() == turnIdx {
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
