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
			inPon:      [3]*hai.Hai{hai.Haku, hai.Haku, hai.Haku},
			afterPons:  [][3]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{pons: c.beforePons}
		err := h.SetPon(c.inPon)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterPons, h.pons)
	}
}

func TestChii(t *testing.T) {
	cases := []struct {
		beforeChiis [][3]*hai.Hai
		inChii      [3]*hai.Hai
		afterChiis  [][3]*hai.Hai
		outError    error
	}{
		{
			beforeChiis: [][3]*hai.Hai{},
			inChii:      [3]*hai.Hai{hai.Haku, hai.Haku, hai.Haku},
			afterChiis:  [][3]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
			outError:    nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{chiis: c.beforeChiis}
		err := h.SetChii(c.inChii)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterChiis, h.chiis)
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
			inKan:      [4]*hai.Hai{hai.Haku, hai.Haku, hai.Haku},
			afterKans:  [][4]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{minKans: c.beforeKans}
		err := h.SetMinKan(c.inKan)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}
		assert.Equal(t, c.afterKans, h.minKans)
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
			beforePons: [][3]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
			inHai:      hai.Haku,
			afterPons:  [][3]*hai.Hai{},
			afterKans:  [][4]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku, hai.Haku}},
			outError:   nil,
		},
		{
			beforePons: [][3]*hai.Hai{{hai.Haku, hai.Haku, hai.Haku}},
			inHai:      hai.Hatu,
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
		assert.Equal(t, c.afterKans, h.minKans)
	}
}
