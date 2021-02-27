package yama

import (
	"mahjong/model/hai"
	"math/rand"
	"time"
)

type Yama interface {
	SetYamaHai([]*hai.Hai) error
	Draw() (*hai.Hai, error)
	Kan() error

	OmoteDora() []*hai.Hai
	UraDora() []*hai.Hai
}

type yamaImpl struct {
	yamaHai   []*hai.Hai
	wanHai    []*hai.Hai
	uraDora   []*hai.Hai
	omoteDora []*hai.Hai
}

var (
	all = []*hai.Hai{
		hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6, hai.Manzu7, hai.Manzu8, hai.Manzu9,
		hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6, hai.Manzu7, hai.Manzu8, hai.Manzu9,
		hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6, hai.Manzu7, hai.Manzu8, hai.Manzu9,
		hai.Manzu1, hai.Manzu2, hai.Manzu3, hai.Manzu4, hai.Manzu5, hai.Manzu6, hai.Manzu7, hai.Manzu8, hai.Manzu9,

		hai.Pinzu1, hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7, hai.Pinzu8, hai.Pinzu9,
		hai.Pinzu1, hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7, hai.Pinzu8, hai.Pinzu9,
		hai.Pinzu1, hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7, hai.Pinzu8, hai.Pinzu9,
		hai.Pinzu1, hai.Pinzu2, hai.Pinzu3, hai.Pinzu4, hai.Pinzu5, hai.Pinzu6, hai.Pinzu7, hai.Pinzu8, hai.Pinzu9,

		hai.Souzu1, hai.Souzu2, hai.Souzu3, hai.Souzu4, hai.Souzu5, hai.Souzu6, hai.Souzu7, hai.Souzu8, hai.Souzu9,
		hai.Souzu1, hai.Souzu2, hai.Souzu3, hai.Souzu4, hai.Souzu5, hai.Souzu6, hai.Souzu7, hai.Souzu8, hai.Souzu9,
		hai.Souzu1, hai.Souzu2, hai.Souzu3, hai.Souzu4, hai.Souzu5, hai.Souzu6, hai.Souzu7, hai.Souzu8, hai.Souzu9,
		hai.Souzu1, hai.Souzu2, hai.Souzu3, hai.Souzu4, hai.Souzu5, hai.Souzu6, hai.Souzu7, hai.Souzu8, hai.Souzu9,

		hai.Ton, hai.Nan, hai.Sha, hai.Pei,
		hai.Ton, hai.Nan, hai.Sha, hai.Pei,
		hai.Ton, hai.Nan, hai.Sha, hai.Pei,
		hai.Ton, hai.Nan, hai.Sha, hai.Pei,

		hai.Haku, hai.Hatsu, hai.Chun,
		hai.Haku, hai.Hatsu, hai.Chun,
		hai.Haku, hai.Hatsu, hai.Chun,
		hai.Haku, hai.Hatsu, hai.Chun,
	}
)

var (
	TimeNowUnix = time.Now().Unix()
)

func New() Yama {
	allHai := append([]*hai.Hai{}, all...)
	rand.Seed(TimeNowUnix)
	rand.Shuffle(len(allHai), func(i, j int) { allHai[i], allHai[j] = allHai[j], allHai[i] })
	return &yamaImpl{
		yamaHai:   allHai[:122],
		wanHai:    allHai[122:],
		omoteDora: []*hai.Hai{},
		uraDora:   []*hai.Hai{},
	}
}

func (y *yamaImpl) SetYamaHai(hais []*hai.Hai) error {
	y.yamaHai = hais
	return nil
}

func (y *yamaImpl) OmoteDora() []*hai.Hai {
	return y.omoteDora
}

func (y *yamaImpl) UraDora() []*hai.Hai {
	return y.uraDora
}

func (y *yamaImpl) Draw() (*hai.Hai, error) {
	if len(y.yamaHai)+(len(y.wanHai)/3) == 4 {
		return nil, YamaNoMoreHaiErr
	}
	outHai := y.yamaHai[0]
	y.yamaHai = y.yamaHai[1:]

	return outHai, nil
}

func (y *yamaImpl) Kan() error {
	if len(y.wanHai) < 2 {
		return YamaNoMoreHaiErr
	}
	od := y.wanHai[0]
	ud := y.wanHai[1]
	y.omoteDora = append(y.omoteDora, od)
	y.uraDora = append(y.uraDora, ud)

	y.wanHai = y.wanHai[2:]
	return nil
}
