package tehai

import (
	"mahjong/hai"
)

type Tehai interface {
	Add(inHai *hai.Hai) error
	Adds(inHais []*hai.Hai) error
	Remove(outhai *hai.Hai) (*hai.Hai, error)
	Removes(outhais []*hai.Hai) ([]*hai.Hai, error)
	Replace(inHai *hai.Hai, outhai *hai.Hai) (*hai.Hai, error)
	FindChiPair(inHai *hai.Hai) [][2]*hai.Hai
	FindPonPair(inHai *hai.Hai) [][2]*hai.Hai
	FindKanPair(inHai *hai.Hai) [][3]*hai.Hai
}
