package utils

import (
	"mahjong/model/hai"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func DrawTehai(hais []*hai.Hai) string {
	l := len(hais)
	str := ""
	for i := 0; i < l; i++ {
		str += "┌──┐"
	}
	str += "\n"
	for _, h := range hais {
		str += "│" + h.Name() + "│"
	}
	str += "\n"
	for i := 0; i < l; i++ {
		str += "└──┘"
	}
	return str
}

func DrawHo(hais []*hai.Hai) string {
	lines := [][]*hai.Hai{}
	for i := 0; i < (len(hais)-1)/6+1; i++ {
		lines = append(lines, []*hai.Hai{})
	}
	for i, h := range hais {
		lines[i/6] = append(lines[i/6], h)
	}
	return DrawHais(lines)

}

func DrawHais(hais [][]*hai.Hai) string {
	str := ""
	for i, _ := range hais {
		top := ""
		body := ""
		bottom := ""
		for _, h := range hais[i] {
			if i == 0 {
				top += "┌──┐"
			}
			body += "│" + h.Name() + "│"
			bottom += "└──┘"
		}
		if i == len(hais)-1 {
			for j := 0; j < len(hais[i-1])-len(hais[i]); j++ {
				body += "└──┘"
			}
			bottom += "\n"
			for j := 0; j < len(hais[i]); j++ {
				bottom += "└──┘"
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
