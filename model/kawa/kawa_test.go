package kawa

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
			inHai:      hai.Haku,
			afterHais:  []*hai.Hai{hai.Haku},
			outError:   nil,
		},
	}

	for _, c := range cases {
		h := kawaImpl{hais: c.beforeHais}
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
			beforeHais: []*hai.Hai{hai.Haku},
			outHai:     hai.Haku,
			outError:   nil,
		},
		{
			beforeHais: []*hai.Hai{},
			outError:   KawaNoHaiError,
		},
	}

	for _, c := range cases {
		h := kawaImpl{hais: c.beforeHais}
		outHai, err := h.Last()

		if err != nil {
			assert.Equal(t, c.outError, err)
			continue
		}

		assert.Equal(t, c.outHai, outHai)

	}
}

func TestRemoveLast(t *testing.T) {
	cases := []struct {
		name       string
		beforeHais []*hai.Hai
		outHai     *hai.Hai
		outError   error
		afterHais  []*hai.Hai
	}{
		{
			name:       "success",
			beforeHais: []*hai.Hai{hai.Haku},
			outHai:     hai.Haku,
			afterHais:  []*hai.Hai{},
		},
		{
			name:       "failure",
			beforeHais: []*hai.Hai{},
			outError:   KawaNoHaiError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Kawa := kawaImpl{hais: c.beforeHais}
			hai, err := Kawa.RemoveLast()
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHai, hai)
			assert.Equal(t, c.afterHais, Kawa.hais)

		})

	}

}
