package cha

import "mahjong/hai"

type Cha interface {
	Tumo() error
	Dahai(outHai *hai.Hai) (*hai.Hai, error)
}
