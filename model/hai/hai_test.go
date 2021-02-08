package hai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtoHai(t *testing.T) {
	cases := []struct {
		name      string
		inHaiName string
		outHai    *Hai
		outError  error
	}{
		{
			name:      "success: manzu",
			inHaiName: "m1",
			outHai:    Manzu1,
		},
		{
			name:      "success: souzu",
			inHaiName: "s2",
			outHai:    Souzu2,
		},
		{
			name:      "success: pinzu",
			inHaiName: "p3",
			outHai:    Pinzu3,
		},
		{
			name:      "success: jihai",
			inHaiName: "東",
			outHai:    Ton,
		},
		{
			name:      "failure",
			inHaiName: "xxx",
			outError:  HaiInvalidArgumentErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hai, err := AtoHai(c.inHaiName)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHai, hai)
		})
	}
}

func TestHaitoI(t *testing.T) {
	cases := []struct {
		name     string
		inHai    *Hai
		outInt   int
		outError error
	}{
		{
			name:   "success: one",
			inHai:  Manzu1,
			outInt: 1,
		},
		{
			name:   "success: two",
			inHai:  Souzu2,
			outInt: 2,
		},
		{
			name:   "success: three",
			inHai:  Pinzu3,
			outInt: 3,
		},
		{
			name:     "failure: jihai",
			inHai:    Ton,
			outError: HaiInvalidArgumentErr,
		},
		{
			name:     "failure: nil",
			inHai:    nil,
			outError: HaiInvalidArgumentErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			num, err := HaitoI(c.inHai)
			if err != nil {
				assert.Equal(t, c.outError, err)

			}
			assert.Equal(t, c.outInt, num)
		})
	}
}

func TestHaitoSuits(t *testing.T) {
	cases := []struct {
		name     string
		inHai    *Hai
		outHais  []*Hai
		outError error
	}{
		{
			name:    "success: manzu",
			inHai:   Manzu1,
			outHais: Manzu,
		},
		{
			name:    "success: souzu",
			inHai:   Souzu2,
			outHais: Souzu,
		},
		{
			name:    "success: pinzu",
			inHai:   Pinzu3,
			outHais: Pinzu,
		},
		{
			name:     "failure: jihai",
			inHai:    Ton,
			outError: HaiInvalidArgumentErr,
		},
		{
			name:     "failure: nil",
			inHai:    nil,
			outError: HaiInvalidArgumentErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hais, err := HaitoSuits(c.inHai)
			if err != nil {
				assert.Equal(t, c.outError, err)
				return
			}
			assert.Equal(t, c.outHais, hais)
		})
	}
}
