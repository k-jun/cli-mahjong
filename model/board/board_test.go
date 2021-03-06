package board

import (
	"errors"
	"mahjong/model/hai"
	"mahjong/model/kawa"
	"mahjong/model/player"
	"mahjong/model/tehai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPlayer(t *testing.T) {
	cases := []struct {
		beforePlayers          []*boardPlayer
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inPlayer               player.Player
		afterPlayersLen        int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforePlayers:          []*boardPlayer{},
			beforeMaxNumberOfUsers: 1,
			beforeIsPlaying:        true,
			inPlayer:               &player.PlayerMock{},
			afterPlayersLen:        1,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforePlayers:          []*boardPlayer{{Player: &player.PlayerMock{}, channel: make(chan Board)}, {Player: &player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               &player.PlayerMock{},
			afterPlayersLen:        3,
			afterIsPlaying:         true,
			outError:               nil,
		},
		{
			beforePlayers:          []*boardPlayer{{Player: &player.PlayerMock{}, channel: make(chan Board)}, {Player: &player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               &player.PlayerMock{},
			outError:               BoardMaxNOUErr,
		},
	}

	for _, c := range cases {
		tk := &boardImpl{
			players:         c.beforePlayers,
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
		assert.Equal(t, c.afterPlayersLen, len(tk.players))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestLeavePlayer(t *testing.T) {

	testPlayer := &player.PlayerMock{}
	cases := []struct {
		beforePlayers          []*boardPlayer
		beforeMaxNumberOfUsers int
		beforeIsPlaying        bool
		inPlayer               player.Player
		afterPlayersLen        int
		afterIsPlaying         bool
		outError               error
	}{
		{
			beforePlayers:          []*boardPlayer{{Player: testPlayer, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*boardPlayer{{Player: testPlayer, channel: make(chan Board)}, {Player: &player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*boardPlayer{{Player: testPlayer, channel: make(chan Board)}, {Player: &player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 3,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			afterPlayersLen:        0,
			afterIsPlaying:         false,
			outError:               nil,
		},
		{
			beforePlayers:          []*boardPlayer{{Player: &player.PlayerMock{}, channel: make(chan Board)}, {Player: &player.PlayerMock{}, channel: make(chan Board)}},
			beforeMaxNumberOfUsers: 2,
			beforeIsPlaying:        true,
			inPlayer:               testPlayer,
			outError:               BoardPlayerNotFoundErr,
		},
	}

	for _, c := range cases {
		tk := &boardImpl{
			players:         c.beforePlayers,
			isPlaying:       c.beforeIsPlaying,
			maxNumberOfUser: c.beforeMaxNumberOfUsers,
		}
		err := tk.LeavePlayer(c.inPlayer)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		if c.beforeIsPlaying != c.afterIsPlaying {
			for _, Player := range tk.players {
				if Board := <-Player.channel; Board != nil {
					t.Fatal()
				}
			}
		} else {
			for _, Player := range tk.players {
				tk2 := <-Player.channel
				assert.Equal(t, tk, tk2)
			}
		}

		assert.Equal(t, c.afterPlayersLen, len(tk.players))
		assert.Equal(t, c.afterIsPlaying, tk.isPlaying)
	}
}

func TestMyTurn(t *testing.T) {
	testPlayer1 := &player.PlayerMock{}
	testPlayer2 := &player.PlayerMock{}
	cases := []struct {
		name          string
		beforePlayers []*boardPlayer
		inPlayer      player.Player
		outInt        int
		outError      error
	}{
		{
			name:          "success",
			beforePlayers: []*boardPlayer{&boardPlayer{Player: testPlayer1}, &boardPlayer{Player: testPlayer2}},
			inPlayer:      testPlayer1,
			outInt:        0,
		},
		{
			name:          "failure",
			beforePlayers: []*boardPlayer{&boardPlayer{Player: testPlayer1}, &boardPlayer{Player: testPlayer1}},
			inPlayer:      testPlayer2,
			outError:      BoardPlayerNotFoundErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := boardImpl{
				players: c.beforePlayers,
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
	testPlayer1 := &player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	TehaiMock1 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{}}
	TehaiMock2 := &tehai.TehaiMock{ChiiMock: [][2]*hai.Hai{{}}}
	testPlayer2 := &player.PlayerMock{TehaiMock: TehaiMock1}
	testPlayer3 := &player.PlayerMock{TehaiMock: TehaiMock2}
	cases := []struct {
		name              string
		beforePlayers     []*boardPlayer
		beforeTurnIndex   int
		afterActionPlayer []*boardPlayer
		outError          error
	}{
		{
			name:              "success: no actions",
			beforePlayers:     []*boardPlayer{{Player: testPlayer1}, {Player: testPlayer2}},
			beforeTurnIndex:   0,
			afterActionPlayer: []*boardPlayer{},
		},
		{
			name:              "success: actions",
			beforePlayers:     []*boardPlayer{{Player: testPlayer1}, {Player: testPlayer3}},
			beforeTurnIndex:   0,
			afterActionPlayer: []*boardPlayer{{Player: testPlayer3}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := boardImpl{
				players:         c.beforePlayers,
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
	testPlayer1 := &player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	testPlayer2 := &player.PlayerMock{KawaMock: &kawa.KawaMock{ErrorMock: errors.New("")}}
	cases := []struct {
		name            string
		beforePlayers   []*boardPlayer
		beforeTurnIndex int
		outHai          *hai.Hai
		outError        error
	}{
		{
			name:            "success",
			beforePlayers:   []*boardPlayer{{Player: testPlayer1}},
			beforeTurnIndex: 0,
			outHai:          hai.Haku,
		},
		{
			name:            "failure",
			beforePlayers:   []*boardPlayer{{Player: testPlayer2}},
			beforeTurnIndex: 0,
			outError:        errors.New(""),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := boardImpl{
				players:   c.beforePlayers,
				turnIndex: c.beforeTurnIndex,
			}
			hai, err := Board.LastKawa()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHai, hai)
		})
	}
}

func TestCancelAction(t *testing.T) {
	testPlayer1 := &player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	testPlayer2 := &player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	cases := []struct {
		name                string
		inPlayer            player.Player
		beforeActionPlayers []*boardPlayer
		outError            error
		afterActionPlayer   []*boardPlayer
	}{
		{
			name:                "success: 2 actioner",
			inPlayer:            testPlayer1,
			beforeActionPlayers: []*boardPlayer{{Player: testPlayer1}, {Player: testPlayer2}},
			afterActionPlayer:   []*boardPlayer{{Player: testPlayer2}},
		},
		{
			name:                "success: after action taken",
			beforeActionPlayers: []*boardPlayer{},
			afterActionPlayer:   []*boardPlayer{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := boardImpl{
				actionPlayers:   c.beforeActionPlayers,
				maxNumberOfUser: 1,
				players:         []*boardPlayer{{}, {}},
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
	testPlayer1 := &player.PlayerMock{KawaMock: &kawa.KawaMock{HaiMock: hai.Haku}}
	cases := []struct {
		name                string
		beforeActionPlayers []*boardPlayer
		beforePlayers       []*boardPlayer
		inPlayer            player.Player
		inFunc              func(*hai.Hai) error
		outError            error
		afterActionPlayer   []*boardPlayer
	}{
		{
			name:                "success",
			beforeActionPlayers: []*boardPlayer{{Player: testPlayer1}},
			beforePlayers:       []*boardPlayer{{Player: testPlayer1}},
			inPlayer:            testPlayer1,
			inFunc:              func(_ *hai.Hai) error { return nil },

			afterActionPlayer: []*boardPlayer{},
		},
		{
			name:                "failure",
			beforePlayers:       []*boardPlayer{},
			beforeActionPlayers: []*boardPlayer{},
			inFunc:              func(_ *hai.Hai) error { return nil },
			outError:            BoardActionAlreadyTokenErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Board := boardImpl{
				actionPlayers: c.beforeActionPlayers,
				players:       c.beforePlayers,
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
