package huro

import (
	"mahjong/hai"
)

type Huro interface {
	AddSet([]*hai.Hai) error
	AddHaiToSet(inTile *hai.Hai) error
	Sets() [][]*hai.Hai
}
