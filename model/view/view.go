package view

import (
	"mahjong/model/board"
	"mahjong/model/hai"
	"mahjong/model/player"
	"strings"
)

type boardViewHai struct {
	*hai.Hai
	isOpen   bool
	isDown   bool
	isRiichi bool
}

type boardViewPlayer struct {
	hais [20]*boardViewHai
}

type boardViewBoard struct {
	hais [20][20]*boardViewHai
}

func NewBoardPlayer(p player.Player, isOpen bool) *boardViewPlayer {
	hais := [20]*boardViewHai{}
	for i, h := range p.Tehai().Hais() {
		hais[i] = &boardViewHai{Hai: h, isOpen: isOpen}
	}
	if p.Tsumohai() != nil {
		hais[len(p.Tehai().Hais())+1] = &boardViewHai{Hai: p.Tsumohai(), isOpen: isOpen}
	}

	head := 20
	// chii
	for _, meld := range p.Naki().Chiis() {
		for _, h := range meld {
			head--
			hais[head] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// pon
	for _, meld := range p.Naki().Pons() {
		for _, h := range meld {
			head--
			hais[head] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
		}
	}
	// ankan
	for _, meld := range p.Naki().AnKans() {
		for i, h := range meld {
			isOpen := true
			if i == 1 || i == 2 {
				isOpen = false
			}
			head--
			hais[head] = &boardViewHai{Hai: h, isOpen: isOpen, isDown: true}
		}
	}
	// minkan
	for _, meld := range p.Naki().MinKans() {
		for _, h := range meld {
			head--
			hais[head] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
		}
	}

	// riichi
	if p.IsRiichi() {
		for _, h := range hais {
			if h != nil {
				h.isRiichi = true
			}
		}
	}
	return &boardViewPlayer{hais}

}

func TehaiOpen(p player.Player) *boardViewPlayer {
	return NewBoardPlayer(p, true)
}

func TehaiHide(p player.Player) *boardViewPlayer {
	return NewBoardPlayer(p, false)
}

func (p *boardViewPlayer) String() string {
	strs := []string{"", "", "", ""}
	for _, h := range p.hais {

		if h == nil {
			strs[0] += "    "
			strs[1] += "    "
			strs[2] += "    "
			strs[3] += "    "
			continue
		}

		lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
		if h != nil && h.isRiichi {
			lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
		}
		if h.isDown {
			if h.isOpen {
				strs[0] += lines[0] + lines[1] + lines[1] + lines[2]
				strs[1] += lines[3] + h.Name() + lines[5]
				strs[2] += lines[6] + lines[7] + lines[7] + lines[8]
				strs[3] += lines[6] + lines[7] + lines[7] + lines[8]
			} else {
				strs[0] += lines[0] + lines[1] + lines[1] + lines[2]
				strs[1] += lines[3] + lines[4] + lines[4] + lines[5]
				strs[2] += lines[6] + lines[7] + lines[7] + lines[8]
				strs[3] += lines[6] + lines[7] + lines[7] + lines[8]
			}
		} else {
			strs[0] += "    "
			strs[1] += lines[0] + lines[1] + lines[1] + lines[2]
			strs[2] += lines[3] + h.Name() + lines[5]
			strs[3] += lines[6] + lines[7] + lines[7] + lines[8]
		}
	}
	return strings.Join(strs, "\n") + "\n"
}

func TehaiKamichaShimochaAndKawaAll(p player.Player, b board.Board) *boardViewBoard {
	hais := [20][20]*boardViewHai{}
	idx, _ := b.MyTurn(p)

	players := b.Players()
	myself := players[idx]
	shimocha := players[(idx+1)%b.MaxNumberOfUser()]
	toimen := players[(idx+2)%b.MaxNumberOfUser()]
	kamicha := players[(idx+3)%b.MaxNumberOfUser()]

	// tehai
	for i, h := range TehaiHide(kamicha).hais {
		if h == nil {
			continue
		}
		hais[i][0] = h
	}
	for i, h := range TehaiHide(shimocha).hais {
		if h == nil {
			continue
		}
		hais[len(hais)-1-i][len(hais[i])-1] = h
	}

	// kawa
	for i, h := range myself.Kawa().Hais() {
		hais[13+i/6][7+i%6] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range shimocha.Kawa().Hais() {
		hais[12-i%6][13+i/6] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range toimen.Kawa().Hais() {
		hais[6-i/6][12-i%6] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
	}
	for i, h := range kamicha.Kawa().Hais() {
		hais[7+i%6][6-i/6] = &boardViewHai{Hai: h, isOpen: true, isDown: true}
	}

	return &boardViewBoard{hais}

}

func (b *boardViewBoard) String() string {
	str := ""
	for i, _ := range b.hais {
		body := ""
		bottom := ""
		top := ""

		for j, h := range b.hais[i] {
			lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
			if (h != nil && h.isRiichi) ||
				(i != 0 && b.hais[i-1][j] != nil && b.hais[i-1][j].isRiichi) ||
				(i != len(b.hais)-1 && b.hais[i+1][j] != nil && b.hais[i+1][j].isRiichi) {
				lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
			}
			if i == 0 {
				if h != nil {
					top += lines[0] + lines[1] + lines[1] + lines[2]
				} else {
					top += "    "
				}
			}
			if h == nil {
				if i != 0 && b.hais[i-1][j] != nil {
					body += lines[6] + lines[7] + lines[7] + lines[8]
				} else {
					body += "    "
				}
				if i != len(b.hais)-1 && b.hais[i+1][j] != nil {
					bottom += lines[0] + lines[1] + lines[1] + lines[2]
				} else {
					bottom += "    "
				}
				continue
			}
			if h.isOpen && h.isDown {
				body += lines[3] + h.Name() + lines[5]
			} else {
				body += lines[3] + lines[4] + lines[4] + lines[5]
			}
			bottom += lines[6] + lines[7] + lines[7] + lines[8]
		}

		if i == len(b.hais)-1 {
			bottom += "\n"
			for _, h := range b.hais[i] {
				lines := []string{"┌", "─", "┐", "│", " ", "│", "└", "─", "┘"}
				if h != nil && h.isRiichi {
					lines = []string{"┏", "━", "┓", "┃", " ", "┃", "┗", "━", "┛"}
				}
				if h == nil {
					bottom += "    "
				} else {
					bottom += lines[6] + lines[7] + lines[7] + lines[8]
				}
			}
		}

		if top != "" {
			str += top + "\n"
		}
		str += body + "\n"
		str += bottom + "\n"
	}
	return str
}

func BoardString(p player.Player, b board.Board) (string, error) {
	str := ""
	idx, err := b.MyTurn(p)
	if err != nil {
		return str, err
	}
	toimen := b.Players()[(idx+2)%b.MaxNumberOfUser()]
	str += TehaiHide(toimen).String()
	str += TehaiKamichaShimochaAndKawaAll(p, b).String()
	str += TehaiOpen(p).String()
	return str, nil
}
