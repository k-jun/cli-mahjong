package player

import (
	"errors"
	"mahjong/model/hai"
	"mahjong/model/kawa"
	"mahjong/model/naki"
	"mahjong/model/tehai"
	"mahjong/model/yama"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTsumo(t *testing.T) {
	cases := []struct {
		beforeTsumohai *hai.Hai
		beforeYama     yama.Yama
		afterTsumohai  *hai.Hai
		outError       error
	}{
		{
			beforeTsumohai: nil,
			beforeYama:     &yama.YamaMock{HaiMock: hai.Haku},
			afterTsumohai:  hai.Haku,
			outError:       nil,
		},
		{
			beforeTsumohai: hai.Haku,
			outError:       PlayerAlreadyHaveTsumohaiErr,
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tsumohai: c.beforeTsumohai,
			yama:     c.beforeYama,
		}

		err := Player.Tsumo()
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTsumohai, Player.tsumohai)
	}
}

func TestDahai(t *testing.T) {
	cases := []struct {
		beforeTsumohai *hai.Hai
		beforeTehai    tehai.Tehai
		beforeKawa     kawa.Kawa
		inHai          *hai.Hai
		afterTsumohai  *hai.Hai
		afterTehai     tehai.Tehai
		afterKawa      kawa.Kawa
		outError       error
	}{
		{
			beforeTsumohai: hai.Haku,
			beforeTehai:    &tehai.TehaiMock{},
			beforeKawa:     &kawa.KawaMock{},
			inHai:          hai.Haku,
			afterTsumohai:  nil,
			afterTehai:     &tehai.TehaiMock{},
			afterKawa:      &kawa.KawaMock{HaiMock: hai.Haku},
		},
		{
			beforeTsumohai: hai.Haku,
			beforeTehai:    &tehai.TehaiMock{HaiMock: hai.Manzu1},
			beforeKawa:     &kawa.KawaMock{},
			inHai:          hai.Manzu1,
			afterTsumohai:  nil,
			afterTehai:     &tehai.TehaiMock{HaiMock: hai.Haku},
			afterKawa:      &kawa.KawaMock{HaiMock: hai.Manzu1},
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tsumohai: c.beforeTsumohai,
			tehai:    c.beforeTehai,
			kawa:     c.beforeKawa,
		}

		err := Player.Dahai(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTsumohai, Player.tsumohai)
		assert.Equal(t, c.afterTehai, Player.tehai)
		assert.Equal(t, c.afterKawa, Player.kawa)
	}
}

func TestChi(t *testing.T) {
	cases := []struct {
		beforeTehai    tehai.Tehai
		beforeNaki     naki.Naki
		beforeTsumohai *hai.Hai
		inHai          *hai.Hai
		inHais         [2]*hai.Hai
		afterNaki      naki.Naki
		afterTsumohai  *hai.Hai
		outError       error
	}{
		{
			beforeTehai:    &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Manzu1, hai.Manzu2}},
			beforeNaki:     &naki.NakiMock{},
			beforeTsumohai: hai.Haku,
			inHai:          hai.Manzu3,
			inHais:         [2]*hai.Hai{hai.Manzu1, hai.Manzu2},
			afterNaki:      &naki.NakiMock{ChiiMock: [3]*hai.Hai{hai.Manzu3, hai.Manzu1, hai.Manzu2}},
			afterTsumohai:  nil,
			outError:       nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeNaki:  &naki.NakiMock{},
			inHai:       hai.Haku,
			inHais:      [2]*hai.Hai{hai.Haku, hai.Haku},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Manzu1, hai.Manzu2}},
			beforeNaki:  &naki.NakiMock{ErrorMock: errors.New("")},
			inHai:       hai.Manzu3,
			inHais:      [2]*hai.Hai{hai.Manzu1, hai.Manzu2},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tehai: c.beforeTehai,
			naki:  c.beforeNaki,
		}

		err := Player.Chii(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterNaki, Player.naki)
		assert.Equal(t, c.afterTsumohai, Player.tsumohai)
	}
}

func TestPon(t *testing.T) {
	cases := []struct {
		beforeTehai    tehai.Tehai
		beforeNaki     naki.Naki
		beforeTsumohai *hai.Hai
		inHai          *hai.Hai
		inHais         [2]*hai.Hai
		afterNaki      naki.Naki
		afterTsumohai  *hai.Hai
		outError       error
	}{
		{
			beforeTehai:    &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Haku, hai.Haku}},
			beforeNaki:     &naki.NakiMock{},
			beforeTsumohai: hai.Haku,
			inHai:          hai.Haku,
			inHais:         [2]*hai.Hai{hai.Haku, hai.Haku},
			afterNaki:      &naki.NakiMock{PonMock: [3]*hai.Hai{hai.Haku, hai.Haku, hai.Haku}},
			afterTsumohai:  nil,
			outError:       nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeNaki:  &naki.NakiMock{},
			inHai:       hai.Haku,
			inHais:      [2]*hai.Hai{hai.Haku, hai.Haku},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Haku, hai.Haku}},
			beforeNaki:  &naki.NakiMock{ErrorMock: errors.New("")},
			inHai:       hai.Haku,
			inHais:      [2]*hai.Hai{hai.Haku, hai.Haku},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tehai: c.beforeTehai,
			naki:  c.beforeNaki,
		}

		err := Player.Pon(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterNaki, Player.naki)
		assert.Equal(t, c.afterTsumohai, Player.tsumohai)
	}
}

func TestKan(t *testing.T) {
	cases := []struct {
		beforeTehai    tehai.Tehai
		beforeNaki     naki.Naki
		beforeTsumohai *hai.Hai
		inHai          *hai.Hai
		inHais         [3]*hai.Hai
		afterNaki      naki.Naki
		afterTsumohai  *hai.Hai
		outError       error
	}{
		{
			beforeTehai:    &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Haku, hai.Haku, hai.Haku}},
			beforeNaki:     &naki.NakiMock{},
			beforeTsumohai: hai.Haku,
			inHai:          hai.Haku,
			inHais:         [3]*hai.Hai{},
			afterNaki:      &naki.NakiMock{MinKanMock: [4]*hai.Hai{hai.Haku, hai.Haku, hai.Haku, hai.Haku}},
			afterTsumohai:  nil,
			outError:       nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeNaki:  &naki.NakiMock{},
			inHai:       hai.Haku,
			inHais:      [3]*hai.Hai{},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{hai.Haku, hai.Haku, hai.Haku}},
			beforeNaki:  &naki.NakiMock{ErrorMock: errors.New("")},
			inHai:       hai.Haku,
			inHais:      [3]*hai.Hai{},
			afterNaki:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tehai: c.beforeTehai,
			naki:  c.beforeNaki,
		}

		err := Player.MinKan(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterNaki, Player.naki)
		assert.Equal(t, c.afterTsumohai, Player.tsumohai)
	}
}

func TestKakan(t *testing.T) {
	cases := []struct {
		beforeNaki     naki.Naki
		beforeTsumohai *hai.Hai
		afterNaki      naki.Naki
		afterTomohai   *hai.Hai
		outError       error
	}{
		{
			beforeNaki:     &naki.NakiMock{PonMock: [3]*hai.Hai{hai.Haku, hai.Haku, hai.Haku}},
			beforeTsumohai: hai.Haku,
			afterNaki:      &naki.NakiMock{MinKanMock: [4]*hai.Hai{hai.Haku, hai.Haku, hai.Haku, hai.Haku}},
			afterTomohai:   nil,
			outError:       nil,
		},
		{
			beforeNaki:     &naki.NakiMock{PonMock: [3]*hai.Hai{hai.Haku, hai.Haku, hai.Haku}},
			beforeTsumohai: hai.Haku,
			afterNaki:      &naki.NakiMock{MinKanMock: [4]*hai.Hai{hai.Haku, hai.Haku, hai.Haku, hai.Haku}},
			afterTomohai:   nil,
			outError:       nil,
		},
		{
			beforeNaki: &naki.NakiMock{ErrorMock: errors.New("")},
			afterNaki:  nil,
			outError:   errors.New(""),
		},
	}

	for _, c := range cases {
		Player := playerImpl{
			tsumohai: c.beforeTsumohai,
			naki:     c.beforeNaki,
		}

		err := Player.Kakan()
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTomohai, Player.tsumohai)
		assert.Equal(t, c.afterNaki, Player.naki)
	}
}

func TestSetYama(t *testing.T) {
	cases := []struct {
		beforeYama yama.Yama
		inYama     yama.Yama
		afterYama  yama.Yama
		outError   error
	}{
		{
			beforeYama: nil,
			inYama:     &yama.YamaMock{},
			afterYama:  &yama.YamaMock{},
			outError:   nil,
		},
		{
			beforeYama: &yama.YamaMock{},
			outError:   PlayerAlreadyHaveYamaErr,
		},
	}

	for _, c := range cases {
		Player := playerImpl{yama: c.beforeYama}
		err := Player.SetYama(c.inYama)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterYama, Player.yama)
	}
}

func TestHaihai(t *testing.T) {
	cases := []struct {
		beforeYama  yama.Yama
		beforeTehai tehai.Tehai
		outError    error
	}{
		{
			beforeYama:  &yama.YamaMock{},
			beforeTehai: &tehai.TehaiMock{},
			outError:    nil,
		},
		{
			beforeYama:  &yama.YamaMock{},
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			outError:    errors.New(""),
		},
		{
			beforeYama:  &yama.YamaMock{},
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{{}}},
			outError:    PlayerAlreadyDidHaipaiErr,
		},
	}

	for _, c := range cases {
		Player := playerImpl{yama: c.beforeYama, tehai: c.beforeTehai}
		err := Player.Haipai()
		assert.Equal(t, c.outError, err)
	}
}

func TestCanTsumoAgari(t *testing.T) {
	cases := []struct {
		name           string
		beforeTehai    tehai.Tehai
		beforeTsumohai *hai.Hai
		outBool        bool
		outError       error
	}{
		{
			beforeTehai: &tehai.TehaiMock{
				BoolMock: true,
			},
			beforeTsumohai: hai.Hatsu,
			outBool:        true,
		},
		{
			beforeTehai: &tehai.TehaiMock{
				BoolMock: false,
			},
			beforeTsumohai: hai.Hatsu,
			outBool:        false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Player := playerImpl{tehai: c.beforeTehai, tsumohai: c.beforeTsumohai}
			isTsumo, err := Player.CanTsumoAgari()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return

			}
			assert.Equal(t, c.outBool, isTsumo)
		})
	}

}
