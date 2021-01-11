package cha

import (
	"errors"
	"mahjong/hai"
	"mahjong/ho"
	"mahjong/huro"
	"mahjong/tehai"
	"mahjong/yama"
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
			afterHuro:     &huro.HuroMock{KanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
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
			afterHuro:     &huro.HuroMock{KanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			afterTomohai:  nil,
			outError:      nil,
		},
		{
			beforeHuro:    &huro.HuroMock{PonMock: [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku}},
			beforeTumohai: nil,
			inHai:         &hai.Haku,
			afterHuro:     &huro.HuroMock{KanMock: [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
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
	testYama := yama.New()
	cases := []struct {
		beforeYama yama.Yama
		inYama     yama.Yama
		afterYama  yama.Yama
		outError   error
	}{
		{
			beforeYama: nil,
			inYama:     testYama,
			afterYama:  testYama,
			outError:   nil,
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
