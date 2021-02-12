package taku

import (
	"errors"
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/ho"
	"mahjong/model/tehai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinCha(t *testing.T) {
	cases := []struct {
		beforeChas             []*takuCha
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inCha                  cha.Cha
		afterChasLen           int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforeChas:             []*takuCha{},
			beforeMaxNumberOfUsers: 1,
			beforeIsPlaying:        true,
			inCha:                  &cha.ChaMock{},
			afterChasLen:           1,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{Cha: &cha.ChaMock{}, channel: make(chan Taku)}, {Cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  &cha.ChaMock{},
			afterChasLen:           3,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{Cha: &cha.ChaMock{}, channel: make(chan Taku)}, {Cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inCha:                  &cha.ChaMock{},
			outError:               TakuMaxNOUErr,
		},
	}

	for _, c := range cases {
		tk := &takuImpl{
			chas:            c.beforeChas,
			isPlaying:       c.beforeIsPlaying,
			maxNumberOfUser: c.beforeMaxNumberOfUsers,
		}
		channel, err := tk.JoinCha(c.inCha)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		for _, cha := range c.beforeChas {
			tk2 := <-cha.channel
			assert.Equal(t, tk, tk2)
		}

		tk2 := <-channel

		assert.Equal(t, tk, tk2)
		assert.Equal(t, c.afterChasLen, len(tk.chas))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestLeaveCha(t *testing.T) {

	testCha := &cha.ChaMock{}
	cases := []struct {
		beforeChas             []*takuCha
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inCha                  cha.Cha
		afterChasLen           int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforeChas:             []*takuCha{{Cha: testCha, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{Cha: testCha, channel: make(chan Taku)}, {Cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{Cha: testCha, channel: make(chan Taku)}, {Cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{Cha: &cha.ChaMock{}, channel: make(chan Taku)}, {Cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			outError:               TakuChaNotFoundErr,
		},
	}

	for _, c := range cases {
		tk := &takuImpl{
			chas:            c.beforeChas,
			isPlaying:       c.beforeIsPlaying,
			maxNumberOfUser: c.beforeMaxNumberOfUsers,
		}
		err := tk.LeaveCha(c.inCha)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		if c.beforeIsPlaying != c.afterIsPlaying {
			for _, cha := range tk.chas {
				if taku := <-cha.channel; taku != nil {
					t.Fatal()
				}
			}
		} else {
			for _, cha := range tk.chas {
				tk2 := <-cha.channel
				assert.Equal(t, tk, tk2)
			}
		}

		assert.Equal(t, c.afterChasLen, len(tk.chas))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestMyTurn(t *testing.T) {
	testCha1 := &cha.ChaMock{}
	testCha2 := &cha.ChaMock{}
	cases := []struct {
		name       string
		beforeChas []*takuCha
		inCha      cha.Cha
		outInt     int
		outError   error
	}{
		{
			name:       "success",
			beforeChas: []*takuCha{&takuCha{Cha: testCha1}, &takuCha{Cha: testCha2}},
			inCha:      testCha1,
			outInt:     0,
		},
		{
			name:       "failure",
			beforeChas: []*takuCha{&takuCha{Cha: testCha1}, &takuCha{Cha: testCha1}},
			inCha:      testCha2,
			outError:   TakuChaNotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				chas: c.beforeChas,
			}
			turnInt, err := taku.MyTurn(c.inCha)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outInt, turnInt)
		})
	}
}

func TestTurnEnd(t *testing.T) {
	testCha1 := &cha.ChaMock{HoMock: &ho.HoMock{HaiMock: hai.Haku}}
	TehaiMock1 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{}}
	TehaiMock2 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{{}}}
	testCha2 := &cha.ChaMock{TehaiMock: TehaiMock1}
	testCha3 := &cha.ChaMock{TehaiMock: TehaiMock2}
	cases := []struct {
		name            string
		beforeChas      []*takuCha
		beforeTurnIndex int
		afterActionCha  []*takuCha
		outError        error
	}{
		{
			name:            "success: no actions",
			beforeChas:      []*takuCha{{Cha: testCha1}, {Cha: testCha2}},
			beforeTurnIndex: 0,
			afterActionCha:  []*takuCha{},
		},
		{
			name:            "success: actions",
			beforeChas:      []*takuCha{{Cha: testCha1}, {Cha: testCha3}},
			beforeTurnIndex: 0,
			afterActionCha:  []*takuCha{{Cha: testCha3}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				chas:            c.beforeChas,
				turnIndex:       c.beforeTurnIndex,
				maxNumberOfUser: MaxNumberOfUsers,
			}
			err := taku.TurnEnd()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionCha, taku.actionChas)
		})
	}
}

func TestLastHo(t *testing.T) {
	testCha1 := &cha.ChaMock{HoMock: &ho.HoMock{HaiMock: hai.Haku}}
	testCha2 := &cha.ChaMock{HoMock: &ho.HoMock{ErrorMock: errors.New("")}}
	cases := []struct {
		name            string
		beforeChas      []*takuCha
		beforeTurnIndex int
		outHai          *hai.Hai
		outError        error
	}{
		{
			name:            "success",
			beforeChas:      []*takuCha{{Cha: testCha1}},
			beforeTurnIndex: 0,
			outHai:          hai.Haku,
		},
		{
			name:            "failure",
			beforeChas:      []*takuCha{{Cha: testCha2}},
			beforeTurnIndex: 0,
			outError:        errors.New(""),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				chas:      c.beforeChas,
				turnIndex: c.beforeTurnIndex,
			}
			hai, err := taku.LastHo()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHai, hai)
		})
	}
}

func TestCancelAction(t *testing.T) {
	testCha1 := &cha.ChaMock{HoMock: &ho.HoMock{HaiMock: hai.Haku}}
	cases := []struct {
		name             string
		inCha            cha.Cha
		beforeActionChas []*takuCha
		outError         error
		afterActionCha   []*takuCha
	}{
		{
			name:             "success: before action taken",
			inCha:            testCha1,
			beforeActionChas: []*takuCha{{Cha: testCha1}, {}},
			afterActionCha:   []*takuCha{{}},
		},
		{
			name:             "success: after action taken",
			beforeActionChas: []*takuCha{},
			afterActionCha:   []*takuCha{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				actionChas:      c.beforeActionChas,
				maxNumberOfUser: 1,
				chas:            []*takuCha{{}, {}},
			}
			err := taku.CancelAction(c.inCha)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionCha, taku.actionChas)
		})
	}
}

func TestTakeAction(t *testing.T) {
	testCha1 := &cha.ChaMock{HoMock: &ho.HoMock{HaiMock: hai.Haku}}
	cases := []struct {
		name             string
		beforeActionChas []*takuCha
		beforeChas       []*takuCha
		inCha            cha.Cha
		inFunc           func(*hai.Hai) error
		outError         error
		afterActionCha   []*takuCha
	}{
		{
			name:             "success",
			beforeActionChas: []*takuCha{{Cha: testCha1}},
			beforeChas:       []*takuCha{{Cha: testCha1}},
			inCha:            testCha1,
			inFunc:           func(_ *hai.Hai) error { return nil },

			afterActionCha: []*takuCha{},
		},
		{
			name:             "failure",
			beforeChas:       []*takuCha{},
			beforeActionChas: []*takuCha{},
			inFunc:           func(_ *hai.Hai) error { return nil },
			outError:         TakuActionAlreadyTokenErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				actionChas: c.beforeActionChas,
				chas:       c.beforeChas,
			}
			err := taku.TakeAction(c.inCha, c.inFunc)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionCha, taku.actionChas)
		})
	}
}
