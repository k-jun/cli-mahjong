package taku

import (
	"mahjong/cha"
)

type Taku interface {
	JoinCha(cha.Cha) (chan Taku, error)
}
