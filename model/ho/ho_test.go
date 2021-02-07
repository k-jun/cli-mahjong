package ho

import (
	"mahjong/model/hai"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {

	cases := []struct {
		beforeHais []*hai.Hai
		inHai      *hai.Hai
		afterHais  []*hai.Hai
		outError   error
	}{
		{
			beforeHais: []*hai.Hai{},
			inHai:      &hai.Haku,
			afterHais:  []*hai.Hai{&hai.Haku},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := hoImpl{hais: c.beforeHais}
		err := h.Add(c.inHai)
		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.afterHais, h.hais)
	}
}

func TestLast(t *testing.T) {
	cases := []struct {
		beforeHais []*hai.Hai
		outHai     *hai.Hai
		outError   error
	}{
		{
			beforeHais: []*hai.Hai{&hai.Haku},
			outHai:     &hai.Haku,
			outError:   nil,
		},
		{
			beforeHais: []*hai.Hai{},
			outError:   HoNoHaiError,
		},
	}

	for _, c := range cases {
		h := hoImpl{hais: c.beforeHais}
		outHai, err := h.Last()

		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.outHai, outHai)

	}
}
