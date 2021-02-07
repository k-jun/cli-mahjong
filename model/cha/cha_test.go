package cha

import (
	"errors"
	"mahjong/model/hai"
	"mahjong/model/ho"
	"mahjong/model/huro"
	"mahjong/model/tehai"
	"mahjong/model/yama"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTumo(t *testing.T) {
	cases := []struct {
		beforeTumohai *hai.Hai
		beforeYama    yama.Yama
		afterTumohai  *hai.Hai
		outError      error
	}{
		{
			beforeTumohai: nil,
			beforeYama:    &yama.YamaMock{HaiMock: &hai.Haku},
			afterTumohai:  &hai.Haku,
			outError:      nil,
		},
		{
			beforeTumohai: &hai.Haku,
			outError:      ChaAlreadyHaveTumohaiErr,
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tumohai: c.beforeTumohai,
			yama:    c.beforeYama,
		}

		err := cha.Tumo()
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTumohai, cha.tumohai)
	}
}

func TestDahai(t *testing.T) {
	cases := []struct {
		beforeTumohai *hai.Hai
		beforeTehai   tehai.Tehai
		beforeHo      ho.Ho
		inHai         *hai.Hai
		afterTumohai  *hai.Hai
		afterTehai    tehai.Tehai
		afterHo       ho.Ho
		outError      error
	}{
		{
			beforeTumohai: &hai.Haku,
			beforeTehai:   &tehai.TehaiMock{},
			beforeHo:      &ho.HoMock{},
			inHai:         &hai.Haku,
			afterTumohai:  nil,
			afterTehai:    &tehai.TehaiMock{},
			afterHo:       &ho.HoMock{HaiMock: &hai.Haku},
		},
		{
			beforeTumohai: &hai.Haku,
			beforeTehai:   &tehai.TehaiMock{HaiMock: &hai.Manzu1},
			beforeHo:      &ho.HoMock{},
			inHai:         &hai.Manzu1,
			afterTumohai:  nil,
			afterTehai:    &tehai.TehaiMock{HaiMock: &hai.Haku},
			afterHo:       &ho.HoMock{HaiMock: &hai.Manzu1},
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tumohai: c.beforeTumohai,
			tehai:   c.beforeTehai,
			ho:      c.beforeHo,
		}

		err := cha.Dahai(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTumohai, cha.tumohai)
		assert.Equal(t, c.afterTehai, cha.tehai)
		assert.Equal(t, c.afterHo, cha.ho)
	}
}

func TestChi(t *testing.T) {
	cases := []struct {
		beforeTehai   tehai.Tehai
		beforeHuro    huro.Huro
		beforeTumohai *hai.Hai
		inHai         *hai.Hai
		inHais        [2]*hai.Hai
		afterHuro     huro.Huro
		afterTumohai  *hai.Hai
		outError      error
	}{
		{
			beforeTehai:   &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Manzu1, &hai.Manzu2}},
			beforeHuro:    &huro.HuroMock{},
			beforeTumohai: &hai.Haku,
			inHai:         &hai.Manzu3,
			inHais:        [2]*hai.Hai{&hai.Manzu1, &hai.Manzu2},
			afterHuro:     &huro.HuroMock{ChiMock: [3]*hai.Hai{&hai.Manzu3, &hai.Manzu1, &hai.Manzu2}},
			afterTumohai:  nil,
			outError:      nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeHuro:  &huro.HuroMock{},
			inHai:       &hai.Haku,
			inHais:      [2]*hai.Hai{&hai.Haku, &hai.Haku},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Manzu1, &hai.Manzu2}},
			beforeHuro:  &huro.HuroMock{ErrorMock: errors.New("")},
			inHai:       &hai.Manzu3,
			inHais:      [2]*hai.Hai{&hai.Manzu1, &hai.Manzu2},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tehai: c.beforeTehai,
			huro:  c.beforeHuro,
		}

		err := cha.Chi(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterHuro, cha.huro)
		assert.Equal(t, c.afterTumohai, cha.tumohai)
	}
}

func TestPon(t *testing.T) {
	cases := []struct {
		beforeTehai   tehai.Tehai
		beforeHuro    huro.Huro
		beforeTumohai *hai.Hai
		inHai         *hai.Hai
		inHais        [2]*hai.Hai
		afterHuro     huro.Huro
		afterTumohai  *hai.Hai
		outError      error
	}{
		{
			beforeTehai:   &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Haku, &hai.Haku}},
			beforeHuro:    &huro.HuroMock{},
			beforeTumohai: &hai.Haku,
			inHai:         &hai.Haku,
			inHais:        [2]*hai.Hai{&hai.Haku, &hai.Haku},
			afterHuro:     &huro.HuroMock{PonMock: [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			afterTumohai:  nil,
			outError:      nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeHuro:  &huro.HuroMock{},
			inHai:       &hai.Haku,
			inHais:      [2]*hai.Hai{&hai.Haku, &hai.Haku},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Haku, &hai.Haku}},
			beforeHuro:  &huro.HuroMock{ErrorMock: errors.New("")},
			inHai:       &hai.Haku,
			inHais:      [2]*hai.Hai{&hai.Haku, &hai.Haku},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tehai: c.beforeTehai,
			huro:  c.beforeHuro,
		}

		err := cha.Pon(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterHuro, cha.huro)
		assert.Equal(t, c.afterTumohai, cha.tumohai)
	}
}

func TestKan(t *testing.T) {
	cases := []struct {
		beforeTehai   tehai.Tehai
		beforeHuro    huro.Huro
		beforeTumohai *hai.Hai
		inHai         *hai.Hai
		inHais        [3]*hai.Hai
		afterHuro     huro.Huro
		afterTumohai  *hai.Hai
		outError      error
	}{
		{
			beforeTehai:   &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			beforeHuro:    &huro.HuroMock{},
			beforeTumohai: &hai.Haku,
			inHai:         &hai.Haku,
			inHais:        [3]*hai.Hai{},
			afterHuro:     &huro.HuroMock{MinKanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			afterTumohai:  nil,
			outError:      nil,
		},
		{
			beforeTehai: &tehai.TehaiMock{ErrorMock: errors.New("")},
			beforeHuro:  &huro.HuroMock{},
			inHai:       &hai.Haku,
			inHais:      [3]*hai.Hai{},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
		{
			beforeTehai: &tehai.TehaiMock{HaisMock: []*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			beforeHuro:  &huro.HuroMock{ErrorMock: errors.New("")},
			inHai:       &hai.Haku,
			inHais:      [3]*hai.Hai{},
			afterHuro:   nil,
			outError:    errors.New(""),
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tehai: c.beforeTehai,
			huro:  c.beforeHuro,
		}

		err := cha.Kan(c.inHai, c.inHais)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterHuro, cha.huro)
		assert.Equal(t, c.afterTumohai, cha.tumohai)
	}
}

func TestKakan(t *testing.T) {
	cases := []struct {
		beforeHuro    huro.Huro
		beforeTumohai *hai.Hai
		inHai         *hai.Hai
		afterHuro     huro.Huro
		afterTomohai  *hai.Hai
		outError      error
	}{
		{
			beforeHuro:    &huro.HuroMock{PonMock: [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			beforeTumohai: &hai.Haku,
			inHai:         &hai.Haku,
			afterHuro:     &huro.HuroMock{MinKanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			afterTomohai:  nil,
			outError:      nil,
		},
		{
			beforeHuro:    &huro.HuroMock{PonMock: [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			beforeTumohai: nil,
			inHai:         &hai.Haku,
			afterHuro:     &huro.HuroMock{MinKanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			afterTomohai:  nil,
			outError:      nil,
		},
		{
			beforeHuro: &huro.HuroMock{ErrorMock: errors.New("")},
			inHai:      &hai.Haku,
			afterHuro:  nil,
			outError:   errors.New(""),
		},
	}

	for _, c := range cases {
		cha := chaImpl{
			tumohai: c.beforeTumohai,
			huro:    c.beforeHuro,
		}

		err := cha.Kakan(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterTomohai, cha.tumohai)
		assert.Equal(t, c.afterHuro, cha.huro)
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
			outError:   ChaAlreadyHaveYamaErr,
		},
	}

	for _, c := range cases {
		cha := chaImpl{yama: c.beforeYama}
		err := cha.SetYama(c.inYama)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterYama, cha.yama)
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
			outError:    ChaAlreadyDidHaihaiErr,
		},
	}

	for _, c := range cases {
		cha := chaImpl{yama: c.beforeYama, tehai: c.beforeTehai}
		err := cha.Haihai()
		assert.Equal(t, c.outError, err)
	}
}

func TestCanRichi(t *testing.T) {
	cases := []struct {
		name           string
		beforeHuro     huro.Huro
		beforeTehai    tehai.Tehai
		beforeTsumohai *hai.Hai
		outHais        []*hai.Hai
	}{
		{
			name:       "両面",
			beforeHuro: &huro.HuroMock{},
			beforeTehai: &tehai.TehaiMock{
				HaisMock: []*hai.Hai{
					&hai.Manzu1, &hai.Manzu2, &hai.Manzu3, &hai.Manzu4, &hai.Manzu5, &hai.Manzu6,
					&hai.Manzu7, &hai.Manzu8, &hai.Manzu9, &hai.Haku, &hai.Pinzu2, &hai.Pinzu3,
					&hai.Hatu,
				},
			},
			beforeTsumohai: &hai.Hatu,
			outHais:        []*hai.Hai{&hai.Haku},
		},
		{
			name:       "嵌張",
			beforeHuro: &huro.HuroMock{},
			beforeTehai: &tehai.TehaiMock{
				HaisMock: []*hai.Hai{
					&hai.Manzu1, &hai.Manzu1, &hai.Manzu1, &hai.Manzu4, &hai.Manzu5, &hai.Manzu6,
					&hai.Manzu7, &hai.Manzu8, &hai.Manzu9, &hai.Haku, &hai.Pinzu5, &hai.Pinzu3,
					&hai.Hatu,
				},
			},
			beforeTsumohai: &hai.Hatu,
			outHais:        []*hai.Hai{&hai.Haku},
		},
		{
			name:       "双碰",
			beforeHuro: &huro.HuroMock{},
			beforeTehai: &tehai.TehaiMock{
				HaisMock: []*hai.Hai{
					&hai.Manzu1, &hai.Manzu1, &hai.Manzu1, &hai.Manzu4, &hai.Manzu5, &hai.Manzu6,
					&hai.Manzu7, &hai.Manzu8, &hai.Manzu9, &hai.Haku, &hai.Pinzu3, &hai.Pinzu3,
					&hai.Hatu,
				},
			},
			beforeTsumohai: &hai.Hatu,
			outHais:        []*hai.Hai{&hai.Haku},
		},
		{
			name:       "単騎",
			beforeHuro: &huro.HuroMock{},
			beforeTehai: &tehai.TehaiMock{
				HaisMock: []*hai.Hai{
					&hai.Manzu1, &hai.Manzu1, &hai.Manzu1, &hai.Manzu4, &hai.Manzu5, &hai.Manzu6,
					&hai.Manzu7, &hai.Manzu8, &hai.Manzu9, &hai.Pinzu1, &hai.Pinzu2, &hai.Pinzu3,
					&hai.Hatu,
				},
			},
			beforeTsumohai: &hai.Haku,
			outHais:        []*hai.Hai{&hai.Hatu, &hai.Haku},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cha := chaImpl{tehai: c.beforeTehai, tumohai: c.beforeTsumohai, huro: c.beforeHuro}
			hais := cha.CanRichi()
			assert.Equal(t, c.outHais, hais)
		})
	}
}

// func TestRemoveHai(t *testing.T) {
// 	hais := []*hai.Hai{&hai.Haku, &hai.Hatu, &hai.Tyun}
// 	outHai := &hai.Haku
// 	hais, _ = removeHai(hais, outHai)
// 	fmt.Println(hais)
// }
