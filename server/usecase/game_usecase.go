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

type InputCommand struct {
	actionType  board.ActionType
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

	ic := InputCommand{actionType: board.Normal, actionIndex: 0, hai: h}
	if h != nil {
		return &ic, nil
	}

	actions := str.FindAllString(rawstr, 1)
	if len(actions) != 1 {
		return &ic, GameUsecaseInvalidActionErr
	}

	switch actions[0] {
	case "tsumo":
		ic.actionType = board.Tsumo
	case "riichi":
		ic.actionType = board.Riichi
	case "chii":
		ic.actionType = board.Chii
	case "pon":
		ic.actionType = board.Pon
	case "kan":
		ic.actionType = board.Kan
	case "ron":
		ic.actionType = board.Ron
	case "no":
		ic.actionType = board.Cancel
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
	if err := b.SetWinner(p); err != nil {
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

func (gu *gameUsecaseImpl) AnKan(b board.Board, p player.Player, ic *InputCommand) error {
	pairs, err := p.Tehai().AnKanPairs(p.Tsumohai())
	if err != nil {
		return err
	}

	if len(pairs) == 0 {
		return GameUsecaseInvalidActionErr
	}
	if ic.actionIndex >= len(pairs) || ic.actionIndex < 0 {
		return GameUsecaseInvalidActionErr
	}
	if err := p.AnKan(pairs[ic.actionIndex]); err != nil {
		return err
	}
	return p.Tsumo()
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
	err := b.TakeAction(p, func(inHai *hai.Hai) error {
		ok, err := p.Tehai().CanRon(inHai)
		if err != nil {
			return err
		}
		if !ok {
			return GameUsecaseInvalidActionErr
		}
		if err := b.SetWinner(p); err != nil {
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
			case board.Normal:
				err = gu.Normal(b, p, command)
			case board.Tsumo:
				err = gu.Tsumo(b, p, command)
			case board.Riichi:
				err = gu.Riichi(b, p, command)
			case board.Kan:
				err = gu.AnKan(b, p, command)
			default:
				continue
			}
		} else {
			// not my turn
			if b.NextTurn() == turnIdx && command.actionType == board.Chii {
				err = gu.Chii(b, p, command)
			} else {
				switch command.actionType {
				case board.Pon:
					err = gu.Pon(b, p, command)
				case board.Kan:
					err = gu.MinKan(b, p, command)
				case board.Ron:
					err = gu.Ron(b, p, command)
				case board.Cancel:
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

func (gu *gameUsecaseImpl) ChiiChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	hai, err := b.LastKawa()
	if err != nil {
		return str, err
	}
	pairs, err := p.Tehai().ChiiPairs(hai)
	if err != nil {
		return str, err
	}
	if len(pairs) != 0 {
		str += "\n" + "chii>> "
		for i, h := range pairs {
			str += strconv.Itoa(i) + ": (" + h[0].Name() + " " + h[1].Name() + ") "
		}
	}
	return str, nil
}

func (gu *gameUsecaseImpl) PonChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	hai, err := b.LastKawa()
	if err != nil {
		return str, err
	}
	pairs, err := p.Tehai().PonPairs(hai)
	if err != nil {
		return str, err
	}
	if len(pairs) != 0 {
		str += "\n" + "pon>> "
		for i, h := range pairs {
			str += strconv.Itoa(i) + ": (" + h[0].Name() + " " + h[1].Name() + ") "
		}
	}
	return str, nil
}

func (gu *gameUsecaseImpl) MinKanChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	hai, err := b.LastKawa()
	if err != nil {
		return str, err
	}
	pairs, err := p.Tehai().MinKanPairs(hai)
	if err != nil {
		return str, err
	}
	if len(pairs) != 0 {
		str += "\n" + "kan>> "
		for i, h := range pairs {
			str += strconv.Itoa(i) + ": (" + h[0].Name() + " " + h[1].Name() + " " + h[2].Name() + ") "
		}
	}
	return str, nil
}

func (gu *gameUsecaseImpl) AnKanChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	hais, err := p.Tehai().AnKanPairs(p.Tsumohai())
	if err != nil {
		return str, err
	}
	if len(hais) != 0 {
		str += "\n" + "kan>>"
		for i, h := range hais {
			str += strconv.Itoa(i) + ": (" + h[0].Name() + " " + h[1].Name() + " " + h[2].Name() + ") "
		}
	}
	return str, nil
}

func (gu *gameUsecaseImpl) RiichiChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	hais, err := p.Tehai().RiichiHais(p.Tsumohai())
	if err != nil {
		return str, err
	}
	if len(hais) != 0 {
		str += "\n" + "riichi>>"
		for i, h := range hais {
			str += strconv.Itoa(i) + ": (" + h.Name() + ") "
		}
	}
	return str, nil
}

func (gu *gameUsecaseImpl) TsumoAgariChoice(b board.Board, p player.Player) (string, error) {
	str := ""
	ok, err := p.CanTsumoAgari()
	if err != nil {
		return str, err
	}
	if ok {
		str += "\n" + "tsumo>> "
	}
	return str, nil
}

func (gu *gameUsecaseImpl) OutputController(p player.Player, channel chan board.Board) error {
	for {
		b, ok := <-channel
		if !ok {
			return GameUsecaseBoardChannelClosedErr
		}

		str, err := view.BoardString(p, b)
		if err != nil {
			return err
		}
		turnIdx, err := b.MyTurn(p)
		if err != nil {
			return err
		}

		if b.CurrentTurn() == turnIdx {
			// my turn

			// ankan
			ankan, err := gu.AnKanChoice(b, p)
			if err != nil {
				return err
			}
			str += ankan

			// riichi
			riichi, err := gu.RiichiChoice(b, p)
			if err != nil {
				return err
			}
			str += riichi
			// tsumo agari
			tsumo, err := gu.TsumoAgariChoice(b, p)
			if err != nil {
				return err
			}
			str += tsumo

		} else {
			// not my turn
			// naki
			actions, err := b.MyAction(p)
			if err != nil {
				return err
			}

			for _, action := range actions {
				choice := ""
				switch action {
				case board.Chii:
					choice, err = gu.ChiiChoice(b, p)
				case board.Pon:
					choice, err = gu.PonChoice(b, p)
				case board.Kan:
					choice, err = gu.MinKanChoice(b, p)
				case board.Ron:
					choice = "ron>> "
				}
				if err != nil {
					return err
				}
				str += choice
			}
		}

		str += "\n"

		if len(b.ActionPlayers()) == 0 && b.CurrentTurn() == turnIdx {
			str += ">>"
		}

		if err := gu.write(str); err != nil {
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
