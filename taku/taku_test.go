package taku

import (
	"mahjong/cha"
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
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        false,
			inCha:                  &cha.ChaMock{},
			afterChasLen:           1,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        false,
			inCha:                  &cha.ChaMock{},
			afterChasLen:           3,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        false,
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
			beforeIsPlaying:        false,
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
			afterChasLen:           2,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: testCha, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        false,
			inCha:                  testCha,
			afterChasLen:           1,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforeChas:             []*takuCha{{cha: &cha.ChaMock{}, channel: make(chan Taku)}, {cha: &cha.ChaMock{}, channel: make(chan Taku)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        false,
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

func TestNextTurn(t *testing.T) {
	cases := []struct {
		beforeChas      []*takuCha
		beforeTurnIndex int
		inIndex         int
		afterTurnIndex  int
		outError        error
	}{
		{
			beforeChas:      []*takuCha{{}, {}},
			beforeTurnIndex: 0,
			inIndex:         1,
			afterTurnIndex:  1,
			outError:        nil,
		},
		{
			beforeChas:      []*takuCha{{}, {}},
			beforeTurnIndex: 0,
			inIndex:         2,
			outError:        TakuIndexOutOfRangeErr,
		},
	}

	for _, c := range cases {
		taku := takuImpl{
			chas:      c.beforeChas,
			turnIndex: c.beforeTurnIndex,
		}
		err := taku.NextTurn(c.inIndex)

		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterTurnIndex, taku.turnIndex)
	}

}
