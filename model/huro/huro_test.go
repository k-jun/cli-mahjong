package huro

import (
	"mahjong/model/hai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPon(t *testing.T) {
	cases := []struct {
		beforePons [][3]*hai.Hai
		inPon      [3]*hai.Hai
		afterPons  [][3]*hai.Hai
		outError   error
	}{
		{
			beforePons: [][3]*hai.Hai{},
			inPon:      [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku},
			afterPons:  [][3]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{pons: c.beforePons}
		err := h.Pon(c.inPon)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterPons, h.pons)
	}
}

func TestChi(t *testing.T) {
	cases := []struct {
		beforeChis [][3]*hai.Hai
		inChi      [3]*hai.Hai
		afterChis  [][3]*hai.Hai
		outError   error
	}{
		{
			beforeChis: [][3]*hai.Hai{},
			inChi:      [3]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku},
			afterChis:  [][3]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{chis: c.beforeChis}
		err := h.Chi(c.inChi)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterChis, h.chis)
	}
}

func TestKan(t *testing.T) {
	cases := []struct {
		beforeKans [][4]*hai.Hai
		inKan      [4]*hai.Hai
		afterKans  [][4]*hai.Hai
		outError   error
	}{
		{
			beforeKans: [][4]*hai.Hai{},
			inKan:      [4]*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku},
			afterKans:  [][4]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{kans: c.beforeKans}
		err := h.Kan(c.inKan)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterKans, h.kans)
	}
}

func TestKakan(t *testing.T) {
	cases := []struct {
		beforePons [][3]*hai.Hai
		inHai      *hai.Hai
		afterPons  [][3]*hai.Hai
		afterKans  [][4]*hai.Hai
		outError   error
	}{
		{
			beforePons: [][3]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			inHai:      &hai.Haku,
			afterPons:  [][3]*hai.Hai{},
			afterKans:  [][4]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
		{
			beforePons: [][3]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			inHai:      &hai.Hatu,
			outError:   HuroNotFoundErr,
		},
	}

	for _, c := range cases {
		h := huroImpl{pons: c.beforePons}
		err := h.Kakan(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterPons, h.pons)
		assert.Equal(t, c.afterKans, h.kans)
	}
}
