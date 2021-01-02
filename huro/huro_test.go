package huro

import (
	"mahjong/hai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSet(t *testing.T) {
	cases := []struct {
		beforeSets [][]*hai.Hai
		inSet      []*hai.Hai
		afterSets  [][]*hai.Hai
		outError   error
	}{
		{
			beforeSets: [][]*hai.Hai{},
			inSet:      []*hai.Hai{&hai.Haku, &hai.Haku, &hai.Haku},
			afterSets:  [][]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := huroImpl{sets: c.beforeSets}
		err := h.AddSet(c.inSet)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterSets, h.sets)
	}

}

func TestAddHaiToSet(t *testing.T) {
	cases := []struct {
		beforeSets [][]*hai.Hai
		inHai      *hai.Hai
		afterSets  [][]*hai.Hai
		outError   error
	}{
		{
			beforeSets: [][]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			inHai:      &hai.Haku,
			afterSets:  [][]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku, &hai.Haku}},
			outError:   nil,
		},
		{
			beforeSets: [][]*hai.Hai{{&hai.Haku, &hai.Haku, &hai.Haku}},
			inHai:      &hai.Hatu,
			outError:   HuroNoSetFoundErr,
		},
	}

	for _, c := range cases {
		h := huroImpl{sets: c.beforeSets}
		err := h.AddHaiToSet(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterSets, h.sets)
	}
}
