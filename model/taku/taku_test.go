package taku

import (
	"errors"
	"mahjong/model/cha"
	"mahjong/model/hai"
	"mahjong/model/ho"
	"mahjong/model/huro"
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
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  &cha.ChaMock{},
			afterChasLen:           3,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
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
			beforeChas:             []*takuCha{{cha: testCha, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: testCha, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: testCha, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inCha:                  testCha,
			afterChasLen:           0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
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
	}{
		{
			name:       "success",
			beforeChas: []*takuCha{&takuCha{cha: testCha1}, &takuCha{cha: testCha2}},
			inCha:      testCha1,
			outInt:     0,
		},
		{
			name:       "failure",
			beforeChas: []*takuCha{&takuCha{cha: testCha1}, &takuCha{cha: testCha1}},
			inCha:      testCha2,
			outInt:     -1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{
				chas: c.beforeChas,
			}
			turnInt := taku.MyTurn(c.inCha)
			assert.Equal(t, c.outInt, turnInt)
		})
	}
}

func TestTurnEnd(t *testing.T) {
	testCha1 := &cha.ChaMock{HoMock: &ho.HoMock{HaiMock: hai.Haku}}
	testCha2 := &cha.ChaMock{HuroActionsMock: []huro.HuroAction{}}
	testCha3 := &cha.ChaMock{HuroActionsMock: []huro.HuroAction{huro.Kan}}
	cases := []struct {
		name               string
		beforeChas         []*takuCha
		beforeTurnIndex    int
		afterActionCounter int
		outError           error
	}{
		{
			name:               "success: no actions",
			beforeChas:         []*takuCha{&takuCha{cha: testCha1}, &takuCha{cha: testCha2}},
			beforeTurnIndex:    0,
			afterActionCounter: 0,
		},
		{
			name:               "success: actions",
			beforeChas:         []*takuCha{&takuCha{cha: testCha1}, &takuCha{cha: testCha3}},
			beforeTurnIndex:    0,
			afterActionCounter: 1,
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
			assert.Equal(t, c.afterActionCounter, taku.actionCounter)
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
			beforeChas:      []*takuCha{&takuCha{cha: testCha1}},
			beforeTurnIndex: 0,
			outHai:          hai.Haku,
		},
		{
			name:            "failure",
			beforeChas:      []*takuCha{&takuCha{cha: testCha2}},
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
	cases := []struct {
		name                string
		beforeActionCounter int
		outError            error
		afterActionCounter  int
	}{
		{
			name:                "success",
			beforeActionCounter: 2,
			afterActionCounter:  1,
		},
		{
			name:                "failure",
			beforeActionCounter: 0,
			outError:            TakuActionAlreadyTokenErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{actionCounter: c.beforeActionCounter}
			err := taku.CancelAction()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionCounter, taku.actionCounter)
		})
	}
}

func TestTakeAction(t *testing.T) {
	cases := []struct {
		name                string
		beforeActionCounter int
		outError            error
		afterActionCounter  int
	}{
		{
			name:                "success",
			beforeActionCounter: 2,
			afterActionCounter:  1,
		},
		{
			name:                "failure",
			beforeActionCounter: 0,
			outError:            TakuActionAlreadyTokenErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taku := takuImpl{actionCounter: c.beforeActionCounter}
			err := taku.CancelAction()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionCounter, taku.actionCounter)
		})
	}
}
