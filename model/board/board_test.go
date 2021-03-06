package board

import (
	"errors"
	"mahjong/model/hai"
	"mahjong/model/kawa"
	"mahjong/model/tehai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPlayer(t *testing.T) {
	cases := []struct {
		beforePlayers          []*BoardPlayer
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inPlayer               Player.Player
		afterPlayersLen        int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforePlayers:          []*BoardPlayer{},
			beforeMaxNumberOfUsers: 1,
			beforeIsPlaying:        true,
			inPlayer:               &Player.PlayerMock{},
			afterPlayersLen:        1,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforePlayers:          []*BoardPlayer{{Player: &Player.PlayerMock{}, channel: make(chan Board)}, {Player: &Player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               &Player.PlayerMock{},
			afterPlayersLen:        3,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforePlayers:          []*BoardPlayer{{Player: &Player.PlayerMock{}, channel: make(chan Board)}, {Player: &Player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               &Player.PlayerMock{},
			outError:               BoardMaxNOUErr,
		},
	}

	for _, c := range cases {
		tk := &BoardImpl{
			Players:         c.beforePlayers,
			isPlaying:       c.beforeIsPlaying,
			maxNumberOfUser: c.beforeMaxNumberOfUsers,
		}
		channel, err := tk.JoinPlayer(c.inPlayer)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		for _, Player := range c.beforePlayers {
			tk2 := <-Player.channel
			assert.Equal(t, tk, tk2)
		}

		tk2 := <-channel

		assert.Equal(t, tk, tk2)
		assert.Equal(t, c.afterPlayersLen, len(tk.Players))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestLeavePlayer(t *testing.T) {

	testPlayer := &Player.PlayerMock{}
	cases := []struct {
		beforePlayers          []*BoardPlayer
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inPlayer               Player.Player
		afterPlayersLen        int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforePlayers:          []*BoardPlayer{{Player: testPlayer, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*BoardPlayer{{Player: testPlayer, channel: make(chan Board)}, {Player: &Player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*BoardPlayer{{Player: testPlayer, channel: make(chan Board)}, {Player: &Player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*BoardPlayer{{Player: &Player.PlayerMock{}, channel: make(chan Board)}, {Player: &Player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			outError:               BoardchanotFoundErr,
		},
	}

	for _, c := range cases {
		tk := &BoardImpl{
			Players:         c.beforePlayers,
			isPlaying:       c.beforeIsPlaying,
			maxNumberOfUser: c.beforeMaxNumberOfUsers,
		}
		err := tk.LeavePlayer(c.inPlayer)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		if c.beforeIsPlaying != c.afterIsPlaying {
			for _, Player := range tk.Players {
				if Board := <-Player.channel; Board != nil {
					t.Fatal()
				}
			}
		} else {
			for _, Player := range tk.Players {
				tk2 := <-Player.channel
				assert.Equal(t, tk, tk2)
			}
		}

		assert.Equal(t, c.afterPlayersLen, len(tk.Players))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestMyTurn(t *testing.T) {
	testPlayer1 := &Player.PlayerMock{}
	testPlayer2 := &Player.PlayerMock{}
	cases := []struct {
		name          string
		beforePlayers []*BoardPlayer
		inPlayer      Player.Player
		outInt        int
		outError      error
	}{
		{
			name:          "success",
			beforePlayers: []*BoardPlayer{&BoardPlayer{Player: testPlayer1}, &BoardPlayer{Player: testPlayer2}},
			inPlayer:      testPlayer1,
			outInt:        0,
		},
		{
			name:          "failure",
			beforePlayers: []*BoardPlayer{&BoardPlayer{Player: testPlayer1}, &BoardPlayer{Player: testPlayer1}},
			inPlayer:      testPlayer2,
			outError:      BoardchanotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := BoardImpl{
				Players: c.beforePlayers,
			}
			turnInt, err := Board.MyTurn(c.inPlayer)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outInt, turnInt)
		})
	}
}

func TestTurnEnd(t *testing.T) {
	testPlayer1 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	TehaiMock1 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{}}
	TehaiMock2 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{{}}}
	testPlayer2 := &Player.PlayerMock{TehaiMock: TehaiMock1}
	testPlayer3 := &Player.PlayerMock{TehaiMock: TehaiMock2}
	cases := []struct {
		name              string
		beforePlayers     []*BoardPlayer
		beforeTurnIndex   int
		afterActionPlayer []*BoardPlayer
		outError          error
	}{
		{
			name:              "success: no actions",
			beforePlayers:     []*BoardPlayer{{Player: testPlayer1}, {Player: testPlayer2}},
			beforeTurnIndex:   0,
			afterActionPlayer: []*BoardPlayer{},
		},
		{
			name:              "success: actions",
			beforePlayers:     []*BoardPlayer{{Player: testPlayer1}, {Player: testPlayer3}},
			beforeTurnIndex:   0,
			afterActionPlayer: []*BoardPlayer{{Player: testPlayer3}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := BoardImpl{
				Players:         c.beforePlayers,
				turnIndex:       c.beforeTurnIndex,
				maxNumberOfUser: MaxNumberOfUsers,
			}
			err := Board.TurnEnd()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionPlayer, Board.actionPlayers)
		})
	}
}

func TestLastkawa(t *testing.T) {
	testPlayer1 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	testPlayer2 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{ErrorMock: errors.New("")}}
	cases := []struct {
		name            string
		beforePlayers   []*BoardPlayer
		beforeTurnIndex int
		outHai          *hai.Hai
		outError        error
	}{
		{
			name:            "success",
			beforePlayers:   []*BoardPlayer{{Player: testPlayer1}},
			beforeTurnIndex: 0,
			outHai:          hai.Haku,
		},
		{
			name:            "failure",
			beforePlayers:   []*BoardPlayer{{Player: testPlayer2}},
			beforeTurnIndex: 0,
			outError:        errors.New(""),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := BoardImpl{
				Players:   c.beforePlayers,
				turnIndex: c.beforeTurnIndex,
			}
			hai, err := Board.Lastkawa()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHai, hai)
		})
	}
}

func TestCancelAction(t *testing.T) {
	testPlayer1 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	testPlayer2 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	cases := []struct {
		name                string
		inPlayer            Player.Player
		beforeActionPlayers []*BoardPlayer
		outError            error
		afterActionPlayer   []*BoardPlayer
	}{
		{
			name:                "success: 2 actioner",
			inPlayer:            testPlayer1,
			beforeActionPlayers: []*BoardPlayer{{Player: testPlayer1}, {Player: testPlayer2}},
			afterActionPlayer:   []*BoardPlayer{{Player: testPlayer2}},
		},
		{
			name:                "success: after action taken",
			beforeActionPlayers: []*BoardPlayer{},
			afterActionPlayer:   []*BoardPlayer{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := BoardImpl{
				actionPlayers:   c.beforeActionPlayers,
				maxNumberOfUser: 1,
				Players:         []*BoardPlayer{{}, {}},
			}
			err := Board.CancelAction(c.inPlayer)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionPlayer, Board.actionPlayers)
		})
	}
}

func TestTakeAction(t *testing.T) {
	testPlayer1 := &Player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	cases := []struct {
		name                string
		beforeActionPlayers []*BoardPlayer
		beforePlayers       []*BoardPlayer
		inPlayer            Player.Player
		inFunc              func(*hai.Hai) error
		outError            error
		afterActionPlayer   []*BoardPlayer
	}{
		{
			name:                "success",
			beforeActionPlayers: []*BoardPlayer{{Player: testPlayer1}},
			beforePlayers:       []*BoardPlayer{{Player: testPlayer1}},
			inPlayer:            testPlayer1,
			inFunc:              func(_ *hai.Hai) error { return nil },

			afterActionPlayer: []*BoardPlayer{},
		},
		{
			name:                "failure",
			beforePlayers:       []*BoardPlayer{},
			beforeActionPlayers: []*BoardPlayer{},
			inFunc:              func(_ *hai.Hai) error { return nil },
			outError:            BoardActionAlreadyTokenErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := BoardImpl{
				actionPlayers: c.beforeActionPlayers,
				Players:       c.beforePlayers,
			}
			err := Board.TakeAction(c.inPlayer, c.inFunc)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.afterActionPlayer, Board.actionPlayers)
		})
	}
}
